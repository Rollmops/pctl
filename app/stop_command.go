package app

import (
	"fmt"
	"sync"
	"time"
)

func StopCommand(names []string, filters []string, noWait bool) error {
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
		go processState.StopAsync(noWait, &wg, consoleMessageChannel)
	}

	consoleMessageChannel.PrintRelevant(processes)
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
			// TODO kill if configured
			return err
		}
	}

	c.stopped = true
	return nil
}

func (c *ProcessState) StopAsync(noWait bool, wg *sync.WaitGroup, messageChannel chan *ConsoleMessage) error {
	defer wg.Done()
	if !c.Process.IsRunning() {
		c.stopped = true
		messageChannel <- &ConsoleMessage{fmt.Sprintf("Process %s has already stopped\n", c.Process.Config), c.Process}
		return nil
	}
	for {
		if c.IsReadyToStop() {
			err := c.Stop(noWait)
			if err != nil {
				messageChannel <- &ConsoleMessage{FailedColor("Failed to stop %s (%s)\n", c.Process.Config, err), c.Process}
			} else {
				messageChannel <- &ConsoleMessage{OkColor("Stopped process %s\n", c.Process.Config), c.Process}
			}
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
}
