package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

//TODO implement noWait
func StopCommand(names []string, filters []string, noWait bool) error {
	processes, err := CurrentContext.Config.CollectSyncedProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf("no matching process Config for name specifiers: %s", strings.Join(names, ", "))
	}
	processStateMap := NewProcessStateMap(
		&processes, func(p *Process) []string {
			return p.Config.DependsOnInverse
		})

	var wg sync.WaitGroup
	wg.Add(len(*processStateMap))

	for _, processState := range *processStateMap {
		go processState.StopAsync(noWait, &wg)
	}

	wg.Wait()
	return nil
}

func (c *ProcessState) Stop(noWait bool) error {
	err := c.Process.Stop()
	if err != nil {
		return err
	}

	if !noWait {
		err = c.Process.WaitForStop(5*time.Second, 100*time.Millisecond)
		if err != nil {
			return err
		}
	}
	log.Debugf("Stopped process '%s'", c.Process.Config.Name)
	c.stopped = true
	return nil
}

func (c *ProcessState) StopAsync(noWait bool, wg *sync.WaitGroup) error {
	if !c.Process.IsRunning() {
		c.stopped = true
		fmt.Printf(WarningColor("Process '%s' has already stopped\n", c.Process.Config.Name))
		wg.Done()
		return nil
	}
	for {
		if c.IsReadyToStop() {
			err := c.Stop(noWait)
			if err != nil {
				fmt.Printf(FailedColor("Failed to stop '%s' (%s)\n", c.Process.Config.Name, err))
			} else {
				fmt.Printf(OkColor("Stopped process '%s'\n", c.Process.Config.Name))
			}
			wg.Done()
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
}
