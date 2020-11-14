package app

import (
	"fmt"
	"sync"
	"time"
)

func KillProcesses(processes []*Process) error {
	processStateMap, err := NewProcessGraphMap(
		processes, func(p *Process) []string {
			return p.Config.DependsOnInverse
		})
	if err != nil {
		return err
	}

	processStatusChannel := make(chan *ProcessStateChange)
	var wg sync.WaitGroup
	wg.Add(len(*processStateMap))

	go func() {
		wg.Wait()
		close(processStatusChannel)
	}()

	for _, processState := range *processStateMap {
		go processState.KillAsync(&wg, &processStatusChannel)
	}

	for processStateChange := range processStatusChannel {
		switch processStateChange.State {
		case StateKilled:
			fmt.Printf(WarningColor("Killed process %s\n", processStateChange.Process.Config.String()))
		case StateKillingError:
			fmt.Printf(FailedColor("Error during process kill of %s (%s)\n", processStateChange.Process.Config.String(), processStateChange.Error.Error()))
		}
	}

	return nil
}

func (c *ProcessGraphNode) KillAsync(wg *sync.WaitGroup, processStateChannel *chan *ProcessStateChange) {
	defer func() {
		wg.Done()
		c.stopped = true
	}()
	if !c.Process.IsRunning() {
		c.stopped = true
		*processStateChannel <- &ProcessStateChange{StateNotRunning, nil, c.Process}
		return
	}
	for {
		ready, _ := c.IsReadyToStop()
		if ready {
			*processStateChannel <- &ProcessStateChange{StateKilling, nil, c.Process}
			err := c.Process.Kill()
			if err != nil {
				*processStateChannel <- &ProcessStateChange{StateKillingError, err, c.Process}
			} else {
				*processStateChannel <- &ProcessStateChange{StateKilled, nil, c.Process}
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
