package app

import (
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/sirupsen/logrus"
	"strings"
)

func StopCommand(names []string, filters Filters, noWait bool, kill bool) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}

	processes, err = AddDependencyProcesses(processes, true)
	if err != nil {
		return err
	}

	if CurrentContext.Config.PromptForStop && !prompter.YN(fmt.Sprintf("Do you really want to proceed stopping?"), false) {
		return nil
	}
	_, err = StopProcesses(processes, noWait, kill)
	return err
}

func StopProcesses(processes []*Process, noWait bool, kill bool) ([]*Process, error) {
	var failedProcessStops []string
	var stoppedProcesses []*Process

	for _, process := range processes {
		processName := process.Config.String()
		if process.IsRunning() {
			fmt.Printf("Stopping %s...", processName)
			err := StopProcess(process, noWait, kill)
			if err != nil {
				failedProcessStops = append(failedProcessStops, processName)
				fmt.Printf(FailedColor("Failed (%s)\n", err.Error()))
			} else {
				stoppedProcesses = append(stoppedProcesses, process)
				fmt.Printf(OkColor("Ok\n"))
			}
		} else {
			logrus.Debugf("Process '%s' is not running", processName)
		}
	}

	if len(failedProcessStops) > 0 {
		return stoppedProcesses, fmt.Errorf("failed to stop: %s", strings.Join(failedProcessStops, ","))
	}
	return stoppedProcesses, nil
}

func StopProcess(process *Process, noWait bool, kill bool) error {
	processName := process.Config.String()
	logrus.Debugf("Stopping process %s", processName)
	err := process.Stop()
	if err != nil {
		return err
	}
	if noWait {
		return nil
	}
	stopped, err := process.WaitForStop()
	if err != nil {
		return err
	}
	if stopped {
		return nil
	}
	if kill {
		err := process.Kill()
		if err != nil {
			return err
		}
		return nil
	}
	err = fmt.Errorf("timeout")
	return err
}
