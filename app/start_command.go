package app

import (
	"fmt"
	"github.com/Songmu/prompter"
	"strings"
)

func StartCommand(names []string, filters Filters, comment string) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}

	processes, err = AddDependencyProcesses(processes, false)
	if err != nil {
		return err
	}

	if CurrentContext.Config.PromptForStart && !prompter.YN(fmt.Sprintf("Do you really want to start?"), false) {
		return nil
	}

	return StartProcesses(processes, comment)
}

func AddDependencyProcesses(processes ProcessList, reverse bool) (ProcessList, error) {
	dependencyTree, err := CurrentContext.Config.GetDependencyGraph(reverse)
	if err != nil {
		return nil, err
	}
	var allTopSortNames []string
	for _, process := range processes {
		if IsInStringList(allTopSortNames, process.Config.String()) {
			continue
		}
		topSortNames, err := dependencyTree.TopSort(process.Config.String())
		if err != nil {
			return nil, err
		}
		for _, topSortName := range topSortNames {
			if !IsInStringList(allTopSortNames, topSortName) {
				allTopSortNames = append(allTopSortNames, topSortName)
			}
		}
	}
	return CurrentContext.Config.CollectProcessesByNameSpecifiers(allTopSortNames, Filters{}, false)
}

func StartProcesses(processes []*Process, comment string) error {
	var failedProcessStarts []string

	for _, process := range processes {
		processName := process.Config.String()
		if !process.IsRunning() {
			fmt.Printf("Starting %s...", processName)
			err := StartProcess(process, comment)
			if err != nil {
				failedProcessStarts = append(failedProcessStarts, processName)
				fmt.Printf(FailedColor("Failed (%s)\n", err.Error()))
			} else {
				fmt.Printf(OkColor("Ok\n"))
			}
		}
	}
	if len(failedProcessStarts) > 0 {
		return fmt.Errorf("failed to start: %s", strings.Join(failedProcessStarts, ","))
	}
	return nil
}

func StartProcess(process *Process, comment string) error {
	pid, err := process.Start(comment)
	if err != nil {
		return err
	}

	ready, err := process.WaitForStartup(pid)
	if err != nil {
		return err
	}
	if !ready {
		err = fmt.Errorf("startup timeout")
		return err
	}
	return nil
}
