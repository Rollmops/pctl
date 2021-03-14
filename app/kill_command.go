package app

import (
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/sirupsen/logrus"
)

func KillCommand(names []string, filters Filters) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}

	if !prompter.YN(fmt.Sprintf("Do you really want to proceed killing?"), false) {
		return nil
	}

	for _, process := range processes {
		processName := process.Config.String()
		if process.IsRunning() {
			err = process.Kill()
			if err != nil {
				fmt.Printf(FailedColor("Error during process kill of %s (%s)\n", processName, err.Error()))
			} else {
				fmt.Printf(WarningColor("Killed process %s\n", processName))
			}
		} else {
			logrus.Debugf("Process '%s' is not running", processName)
		}
	}
	return nil
}
