package app

import (
	"fmt"
	"sync"
	"time"
)

func StartCommand(names []string, filters []string, comment string) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}

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
		// TODO handle error
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
		*consoleMessageChannel <- &ConsoleMessage{fmt.Sprintf("Waiting %s after starting %s\n", DurationToString(duration), c.Process.Config), c.Process}
		time.Sleep(duration)
	}

	err = c.Process.WaitForReady()
	if err != nil {
		return err
	}
	c.started = true
	return nil
}

func (c *ProcessState) StartAsync(wg *sync.WaitGroup, comment string, consoleMessageChannel *ConsoleMessageChannel) error {
	defer wg.Done()
	if c.Process.IsRunning() {
		c.started = true
		*consoleMessageChannel <- &ConsoleMessage{WarningColor("Process %s is already running\n", c.Process.Config), c.Process}
		return nil
	}
	for {
		if c.IsReadyToStart() {
			err := c.Start(comment, consoleMessageChannel)
			if err != nil {
				*consoleMessageChannel <- &ConsoleMessage{FailedColor("Failed to start %s (%s)\n", c.Process.Config, err), c.Process}
			} else {
				*consoleMessageChannel <- &ConsoleMessage{OkColor("Started process %s\n", c.Process.Config), c.Process}
			}
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
}
