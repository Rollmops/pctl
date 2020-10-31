package app

import (
	"fmt"
)

func RestartCommand(names []string, filters Filters, comment string, kill bool) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	stoppedProcesses, err := StopProcesses(processes, false, kill)
	if err != nil {
		return err
	}
	err = CurrentContext.Cache.Refresh()
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	if err != nil {
		return err
	}

	return StartProcesses(stoppedProcesses, comment)
}
