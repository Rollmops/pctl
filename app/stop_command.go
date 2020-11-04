package app

import (
	"fmt"
	"github.com/Songmu/prompter"
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

	if CurrentContext.Config.PromptForStop && !prompter.YN(fmt.Sprintf("Do you really want to proceed stopping?"), false) {
		return nil
	}
	_, err = StopProcesses(processes, noWait, kill)
	return err
}

func StopProcesses(processes []*Process, noWait bool, kill bool) ([]*Process, error) {
	processStateMap, err := NewProcessGraphMap(
		processes, func(p *Process) []string {
			return p.Config.DependsOnInverse
		})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	processStatusChannel := make(chan *ProcessStateChange)
	wg.Add(len(*processStateMap))

	go func() {
		wg.Wait()
		close(processStatusChannel)
	}()

	for _, processState := range *processStateMap {
		go processState.StopAsync(noWait, kill, &wg, &processStatusChannel)
	}
	var stoppedProcesses []*Process

	for processChange := range processStatusChannel {
		switch processChange.State {
		case StateStopped:
			stoppedProcesses = append(stoppedProcesses, processChange.Process)
			fmt.Printf(OkColor("Stopped process %s\n", processChange.Process.Config.String()))
		case StateStoppingError:
			fmt.Printf(FailedColor("Error during process stop of %s (%s)\n", processChange.Process.Config.String(), processChange.Error.Error()))
		case StateDependencyStoppingError:
			fmt.Printf(WarningColor("Will not stop %s\n", processChange.Process.Config.String()))
		case StateKilled:
			stoppedProcesses = append(stoppedProcesses, processChange.Process)
			fmt.Printf(WarningColor("Killed process %s\n", processChange.Process.Config.String()))
		}
	}

	return stoppedProcesses, nil
}

func (c *ProcessGraphNode) Stop(noWait bool, kill bool, processStatusChannel *chan *ProcessStateChange) error {
	*processStatusChannel <- &ProcessStateChange{StateStopping, nil, c.Process}
	err := c.Process.Stop()
	if err != nil {
		c.stopErr = &err
		*processStatusChannel <- &ProcessStateChange{StateStoppingError, err, c.Process}
		return err
	}
	if noWait {
		*processStatusChannel <- &ProcessStateChange{StateStopped, nil, c.Process}
		c.stopped = true
		return nil
	}
	stopped, err := c.Process.WaitForStop()
	if err != nil {
		c.stopErr = &err
		*processStatusChannel <- &ProcessStateChange{StateStoppingError, err, c.Process}
		return err
	}
	if stopped {
		*processStatusChannel <- &ProcessStateChange{StateStopped, nil, c.Process}
		c.stopped = true
		return nil
	}
	if kill {
		*processStatusChannel <- &ProcessStateChange{StateKilling, nil, c.Process}
		err := c.Process.Kill()
		if err != nil {
			c.stopErr = &err
			*processStatusChannel <- &ProcessStateChange{StateKillingError, err, c.Process}
			return err
		}
		*processStatusChannel <- &ProcessStateChange{StateKilled, nil, c.Process}
		c.stopped = true
		return nil
	}
	err = fmt.Errorf("timeout")
	*processStatusChannel <- &ProcessStateChange{StateStoppingError, err, c.Process}
	c.stopErr = &err
	return err
}

func (c *ProcessGraphNode) StopAsync(noWait bool, kill bool, wg *sync.WaitGroup, processStatusChannel *chan *ProcessStateChange) {
	defer wg.Done()
	if !c.Process.IsRunning() {
		*processStatusChannel <- &ProcessStateChange{StateNotRunning, nil, c.Process}
		c.stopped = true
		return
	}
	for {
		readyToStop, err := c.IsReadyToStop()

		if err != nil {
			*processStatusChannel <- &ProcessStateChange{StateDependencyStoppingError, nil, c.Process}
			c.stopErr = &err
			return
		}
		if readyToStop {
			_ = c.Stop(noWait, kill, processStatusChannel)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
