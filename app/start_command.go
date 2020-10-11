package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type ProcessReadyCheck struct {
	Process      *process.Process
	dependencies []*ProcessReadyCheck
	started      bool
}

// Is ready to start, when all dependencies are started
func (c *ProcessReadyCheck) IsReadyToStart() bool {
	for _, d := range c.dependencies {
		if !d.started {
			return false
		}
	}
	return true
}

func (c *ProcessReadyCheck) Start() error {
	err := c.Process.Start()
	if err != nil {
		return err
	}

	err = c.Process.WaitForStarted(5*time.Second, 100*time.Millisecond)
	if err != nil {
		return err
	}
	log.Debugf("Started process '%s'", c.Process.Config.Name)
	c.started = true
	return nil
}

func (c *ProcessReadyCheck) AddDependency(d *ProcessReadyCheck) {
	for _, dep := range c.dependencies {
		if d == dep {
			return
		}
	}
	c.dependencies = append(c.dependencies, d)
}

func StartProcessReadyCheck(c *ProcessReadyCheck, wg *sync.WaitGroup, s *persistence.Data, comment string) error {
	entry := s.FindByName(c.Process.Config.Name)
	if entry != nil {
		err := c.Process.SynchronizeWithPid(entry.Pid)
		if err == nil && c.Process.IsRunning() {
			c.started = true
			fmt.Printf(output.OkColor("Process '%s' is already running\n", c.Process.Config.Name))
			wg.Done()
			return nil
		}
	}
	for {
		if c.IsReadyToStart() {
			err := c.Start()
			if err != nil {
				s.RemoveByName(c.Process.Config.Name)
				fmt.Printf(output.FailedColor("Failed to start '%s'\n", c.Process.Config.Name))
			} else {
				newDataEntry, err := persistence.NewDataEntryFromProcess(c.Process)
				newDataEntry.Comment = comment
				newDataEntry.MarkFlag = persistence.MarkedAsStarted
				if err != nil {
					return err
				}
				s.AddOrUpdateEntry(newDataEntry)
				fmt.Printf(output.OkColor("Started process '%s'\n", c.Process.Config.Name))
			}
			wg.Done()
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func StartCommand(names []string, comment string) error {
	processConfigs := CurrentContext.Config.CollectProcessConfigsByNameSpecifiers(names, false)
	if len(processConfigs) == 0 {
		return fmt.Errorf("no matching process config for name specifiers: %s", strings.Join(names, ", "))
	}
	persistentState, err := CurrentContext.PersistenceReader.Read()
	if err != nil {
		return nil
	}

	prc := make(map[string]*ProcessReadyCheck)
	var wg sync.WaitGroup

	for _, p := range processConfigs {
		prc = addToProcessReadyCheckMap(p, prc, &wg)
	}
	for _, v := range prc {
		go StartProcessReadyCheck(v, &wg, persistentState, comment)
	}
	wg.Wait()
	return CurrentContext.PersistenceWriter.Write(persistentState)
}

func addToProcessReadyCheckMap(c *config.ProcessConfig, prc map[string]*ProcessReadyCheck, wg *sync.WaitGroup) map[string]*ProcessReadyCheck {
	if prc[c.Name] == nil {
		wg.Add(1)
		prc[c.Name] = &ProcessReadyCheck{
			Process: &process.Process{Config: c},
			started: false,
		}
	}
	for _, d := range c.DependsOn {
		prc = addToProcessReadyCheckMap(CurrentContext.Config.FindByName(d), prc, wg)
		prc[c.Name].AddDependency(prc[d])
	}
	return prc
}

/*
	- get persistence data entry for name
	  - if not present (assume not running), start process
	  - if present, check state
	    - state: running -> do nothing (already running)
		- state: stopped unexpected -> start process
*/
