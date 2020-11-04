package app

import (
	"fmt"
	"github.com/Songmu/prompter"
)

func RestartCommand(names []string, filters Filters, comment string, kill bool) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}
	if CurrentContext.Config.PromptForStop && !prompter.YN(fmt.Sprintf("Do you really want to proceed with restart?"), false) {
		return nil
	}
	stoppedProcesses, err := StopProcesses(processes, false, kill)
	if err != nil {
		return err
	}
	if len(stoppedProcesses) == 0 {
		return nil
	}
	err = CurrentContext.Cache.Refresh()
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	if err != nil {
		return err
	}

	return StartProcesses(stoppedProcesses, comment)
}
