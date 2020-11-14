package app

import (
	"fmt"
	"sync"
	"time"
)

func StartProcesses(processes []*Process, comment string) error {
	processStateMap, err := NewProcessGraphMap(
		processes, func(p *Process) []string {
			return p.Config.DependsOn
		})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	processStatusChannel := make(chan *ProcessStateChange)
	wg.Add(len(*processStateMap))

	go func() {
		wg.Wait()
		close(processStatusChannel)
	}()

	for _, processState := range *processStateMap {
		go processState.StartAsync(&wg, comment, &processStatusChannel)
	}

	for processStateChange := range processStatusChannel {
		switch processStateChange.State {
		case StateStarted:
			fmt.Printf(OkColor("Started process %s\n", processStateChange.Process.Config.String()))
		case StateStartingError:
			fmt.Printf(FailedColor("Error during process start of %s (%s)\n", processStateChange.Process.Config.String(), processStateChange.Error.Error()))
		case StateRunning:
			for _, process := range processes {
				if process.Config.String() == processStateChange.Process.Config.String() {
					fmt.Printf("Process %s is already running", processStateChange.Process.Config.String())
				}
			}
		}
	}
	return nil
}

func (c *ProcessGraphNode) Start(comment string) error {
	pid, err := c.Process.Start(comment)
	if err != nil {
		c.startErr = &err
		return err
	}

	ready, err := c.Process.WaitForStartup(pid)
	if err != nil {
		c.startErr = &err
		return err
	}
	if !ready {
		err = fmt.Errorf("startup timeout")
		return err
	}
	c.started = true

	return nil
}

func (c *ProcessGraphNode) StartAsync(wg *sync.WaitGroup, comment string, processStateChannel *chan *ProcessStateChange) {
	defer wg.Done()
	if c.Process.IsRunning() {
		c.started = true
		*processStateChannel <- &ProcessStateChange{StateRunning, nil, c.Process}
		return
	}
	for {
		ready, err := c.IsReadyToStart()
		if err != nil {
			c.startErr = &err
			return
		}
		if ready {
			*processStateChannel <- &ProcessStateChange{StateStarting, nil, c.Process}
			err := c.Start(comment)
			if err != nil {
				*processStateChannel <- &ProcessStateChange{StateStartingError, err, c.Process}
			} else {
				*processStateChannel <- &ProcessStateChange{StateStarted, nil, c.Process}
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
