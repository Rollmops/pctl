package app

import (
	"fmt"
	"sync"
	"time"
)

func StartCommand(names []string, filters Filters, comment string) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}
	return StartProcesses(processes, comment)

}

func StartProcesses(processes []*Process, comment string) error {
	processStateMap, err := NewProcessStateMap(
		processes, func(p *Process) []string {
			return p.Config.DependsOn
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
		go processState.StartAsync(&wg, comment, &consoleMessageChannel)
	}
	consoleMessageChannel.PrintRelevant(processes)
	return nil
}

func (c *ProcessState) Start(comment string, consoleMessageChannel *ConsoleMessageChannel) error {
	err := c.Process.Start(comment)
	if err != nil {
		return err
	}

	if c.Process.Config.WaitAfterStart != "" {
		duration, err := time.ParseDuration(c.Process.Config.WaitAfterStart)
		if err != nil {
			return err
		}
		*consoleMessageChannel <- &ConsoleMessage{fmt.Sprintf("Waiting %s after start of %s\n", DurationToString(duration), c.Process.Config), nil}
		time.Sleep(duration)
	}

	ready, err := c.Process.WaitForReady()
	if err != nil {
		return err
	}
	if !ready {
		*consoleMessageChannel <- &ConsoleMessage{FailedColor("Unable to start %s\n", c.Process.Config), nil}
		return nil
	}
	c.started = true
	return nil
}

func (c *ProcessState) StartAsync(wg *sync.WaitGroup, comment string, consoleMessageChannel *ConsoleMessageChannel) error {
	defer wg.Done()
	if c.Process.IsRunning() {
		c.started = true
		*consoleMessageChannel <- &ConsoleMessage{fmt.Sprintf("Process %s is already running\n", c.Process.Config), c.Process}
		return nil
	}
	for {
		if c.IsReadyToStart() {
			err := c.Start(comment, consoleMessageChannel)
			if err != nil {
				*consoleMessageChannel <- &ConsoleMessage{FailedColor("Error during starting %s (%s)\n", c.Process.Config, err), nil}
			} else {
				*consoleMessageChannel <- &ConsoleMessage{OkColor("Started process %s\n", c.Process.Config), nil}
			}
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
}
