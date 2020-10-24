package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/process"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func StartCommand(names []string, filters []string, comment string) error {
	processConfigs, err := CurrentContext.Config.CollectProcessConfigsByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processConfigs) == 0 {
		return fmt.Errorf("no matching process config for name specifiers: %s", strings.Join(names, ", "))
	}
	processStateMap := NewFromProcessConfigs(
		&processConfigs, func(c *config.ProcessConfig) []string {
			return c.DependsOn
		})

	var wg sync.WaitGroup
	wg.Add(len(*processStateMap))

	for _, processState := range *processStateMap {
		// TODO handle error
		go processState.StartAsync(&wg, comment)
	}

	wg.Wait()
	return nil
}

func (c *ProcessState) Start(comment string) error {
	err := c.Process.Start(comment)
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

func (c *ProcessState) StartAsync(wg *sync.WaitGroup, comment string) error {
	runningEnvironInfo, err := process.FindRunningInfo(c.Process.Config.Name)
	if err != nil {
		return err
	}
	if runningEnvironInfo != nil {
		c.started = true
		fmt.Printf(output.OkColor("Process '%s' is already running\n", c.Process.Config.Name))
		wg.Done()
		return nil
	}
	for {
		if c.IsReadyToStart() {
			err := c.Start(comment)
			if err != nil {
				fmt.Printf(output.FailedColor("Failed to start '%s' (%s)\n", c.Process.Config.Name, err))
			} else {
				fmt.Printf(output.OkColor("Started process '%s'\n", c.Process.Config.Name))
			}
			wg.Done()
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
}

/*
	- get persistence data entry for name
	  - if not present (assume not running), start process
	  - if present, check state
	    - state: running -> do nothing (already running)
		- state: stopped unexpected -> start process
*/
