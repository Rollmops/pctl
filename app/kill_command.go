package app

import (
	"fmt"
	"sync"
	"time"
)

func KillCommand(names []string, filters []string) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}

	err = CurrentContext.Processes.SyncRunningInfo()
	if err != nil {
		return err
	}

	processStateMap, err := NewProcessStateMap(
		processes, func(p *Process) []string {
			return p.Config.DependsOnInverse
		})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(*processStateMap))

	for _, processState := range *processStateMap {
		go processState.KillAsync(&wg)
	}

	wg.Wait()
	return nil
}

func (c *ProcessState) KillAsync(wg *sync.WaitGroup) error {
	defer func() {
		wg.Done()
		c.stopped = true
	}()
	if !c.Process.IsRunning() {
		c.stopped = true
		fmt.Printf(WarningColor("Process '%s' has already stopped\n", c.Process.Config.Name))
		return nil
	}
	for {
		if c.IsReadyToStop() {
			err := c.Process.Kill()
			if err != nil {
				fmt.Printf(FailedColor("Failed to kill '%s' (%s)\n", c.Process.Config.Name, err))
			} else {
				fmt.Printf(OkColor("Killed process '%s'\n", c.Process.Config.Name))
			}
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
}
