package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func StartCommand(names []string, filters []string, comment string) error {
	processes, err := CurrentContext.Config.CollectSyncedProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf("no matching process config for name specifiers: %s", strings.Join(names, ", "))
	}
	processStateMap := NewProcessStateMap(
		&processes, func(p *Process) []string {
			return p.Config.DependsOn
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
	if c.Process.IsRunning() {
		c.started = true
		fmt.Printf(WarningColor("Process '%s' is already running\n", c.Process.Config.Name))
		wg.Done()
		return nil
	}
	for {
		if c.IsReadyToStart() {
			err := c.Start(comment)
			if err != nil {
				fmt.Printf(FailedColor("Failed to start '%s' (%s)\n", c.Process.Config.Name, err))
			} else {
				fmt.Printf(OkColor("Started process '%s'\n", c.Process.Config.Name))
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
