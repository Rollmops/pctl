package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/process"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

//TODO implement noWait
func StopCommand(names []string, filters []string, noWait bool) error {
	processConfigs, err := CurrentContext.Config.CollectProcessConfigsByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processConfigs) == 0 {
		return fmt.Errorf("no matching process Config for name specifiers: %s", strings.Join(names, ", "))
	}
	processStateMap := NewFromProcessConfigs(
		&processConfigs, func(c *config.ProcessConfig) []string {
			return c.DependsOnInverse
		})

	var wg sync.WaitGroup
	wg.Add(len(*processStateMap))

	for _, processState := range *processStateMap {
		go processState.StopAsync(&wg)
	}

	wg.Wait()
	return nil
}

func (c *ProcessState) Stop() error {
	err := c.Process.Stop()
	if err != nil {
		return err
	}

	err = c.Process.WaitForStop(5*time.Second, 100*time.Millisecond)
	if err != nil {
		return err
	}
	log.Debugf("Stopped process '%s'", c.Process.Config.Name)
	c.stopped = true
	return nil
}

func (c *ProcessState) StopAsync(wg *sync.WaitGroup) error {
	runningEnvironInfo, err := process.FindRunningInfo(c.Process.Config.Name)
	if err != nil {
		return err
	}
	if runningEnvironInfo == nil {
		c.stopped = true
		fmt.Printf(output.WarningColor("Process '%s' has already stopped\n", c.Process.Config.Name))
		wg.Done()
		return nil
	}
	for {
		if c.IsReadyToStop() {
			err := c.Stop()
			if err != nil {
				fmt.Printf(output.FailedColor("Failed to stop '%s' (%s)\n", c.Process.Config.Name, err))
			} else {
				fmt.Printf(output.OkColor("Stopped process '%s'\n", c.Process.Config.Name))
			}
			wg.Done()
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
}
