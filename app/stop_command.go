package app

import (
	"fmt"
	"sync"
	"time"
)

func StopCommand(names []string, filters Filters, noWait bool, kill bool) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}
	_, err = StopProcesses(processes, noWait, kill)
	return err
}

func StopProcesses(processes []*Process, noWait bool, kill bool) ([]*Process, error) {
	processStateMap, err := NewProcessStateMap(
		processes, func(p *Process) []string {
			return p.Config.DependsOnInverse
		})
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	consoleMessageChannel := make(ConsoleMessageChannel)
	wg.Add(len(*processStateMap))

	go func() {
		wg.Wait()
		close(consoleMessageChannel)
	}()

	for _, processState := range *processStateMap {
		go processState.StopAsync(noWait, kill, &wg, &consoleMessageChannel)
	}

	consoleMessageChannel.PrintRelevant(processes)

	var stoppedProcesses []*Process
	for _, v := range *processStateMap {
		stoppedProcesses = append(stoppedProcesses, v.Process)
	}
	return stoppedProcesses, nil
}

func (c *ProcessState) Stop(noWait bool, kill bool, consoleMessageChannel *ConsoleMessageChannel) error {
	err := c.Process.Stop()
	if err != nil {
		c.stopErr = &err
		return err
	}
	if !noWait {
		stopped, err := c.Process.WaitForStop(5*time.Second, 100*time.Millisecond)
		if err != nil {
			c.stopErr = &err
			return err
		}
		if !stopped {
			if kill {
				*consoleMessageChannel <- &ConsoleMessage{WarningColor("Unable to stop %s ... killing\n", c.Process.Config), nil}
				err := c.Process.Kill()
				if err != nil {
					c.stopErr = &err
					return err
				}
				c.stopped = true
			} else {
				err := fmt.Errorf("unable to stop %s", c.Process.Config)
				c.stopErr = &err
				*consoleMessageChannel <- &ConsoleMessage{FailedColor("Unable to stop %s\n", c.Process.Config), nil}
				return nil
			}
		}
	}
	c.stopped = true
	return nil
}

func (c *ProcessState) StopAsync(noWait bool, kill bool, wg *sync.WaitGroup, consoleMessageChannel *ConsoleMessageChannel) error {
	defer wg.Done()
	if !c.Process.IsRunning() {
		c.stopped = true
		*consoleMessageChannel <- &ConsoleMessage{fmt.Sprintf("Process %s has already stopped\n", c.Process.Config), c.Process}
		return nil
	}
	for {
		readyToStop, err := c.IsReadyToStop()
		if err != nil {
			*consoleMessageChannel <- &ConsoleMessage{WarningColor("Will not stop %s\n", c.Process.Config), nil}
			c.stopErr = &err
			return nil
		}
		if readyToStop {
			err := c.Stop(noWait, kill, consoleMessageChannel)
			if err != nil {
				*consoleMessageChannel <- &ConsoleMessage{FailedColor("Error during stopping %s (%s)\n", c.Process.Config, err), nil}
			} else {
				*consoleMessageChannel <- &ConsoleMessage{OkColor("Stopped process %s\n", c.Process.Config), nil}
			}
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
}
