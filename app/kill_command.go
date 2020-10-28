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

	processStateMap, err := NewProcessStateMap(
		processes, func(p *Process) []string {
			return p.Config.DependsOnInverse
		})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	consoleMessageChannel := make(ConsoleMessageChannel)
	wg.Add(len(*processStateMap))

	go func() {
		wg.Wait()
		close(consoleMessageChannel)
	}()

	for _, processState := range *processStateMap {
		go processState.KillAsync(&wg, &consoleMessageChannel)
	}

	consoleMessageChannel.PrintRelevant(processes)
	return nil
}

func (c *ProcessState) KillAsync(wg *sync.WaitGroup, consoleMessageChannel *ConsoleMessageChannel) error {
	defer func() {
		wg.Done()
		c.stopped = true
	}()
	if !c.Process.IsRunning() {
		c.stopped = true
		*consoleMessageChannel <- &ConsoleMessage{WarningColor("Process %s has already stopped\n", c.Process.Config), c.Process}
		return nil
	}
	for {
		if c.IsReadyToStop() {
			err := c.Process.Kill()
			if err != nil {
				*consoleMessageChannel <- &ConsoleMessage{FailedColor("Failed to kill %s (%s)\n", c.Process.Config, err), c.Process}
			} else {
				*consoleMessageChannel <- &ConsoleMessage{OkColor("Killed process %s\n", c.Process.Config), c.Process}
			}
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
}
