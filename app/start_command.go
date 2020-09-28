package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	gopsutil "github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
	"strings"
)

func StartCommand(names []string, comment string) error {
	processConfigs := CurrentContext.config.CollectProcessConfigsByNameSpecifiers(names, false)
	if len(processConfigs) == 0 {
		return fmt.Errorf("no matching process config for name specifiers: %s", strings.Join(names, ", "))
	}
	log.Tracef("Starting processes: %v", names)
	data, err := CurrentContext.persistenceReader.Read()
	if err != nil {
		return nil
	}
	for _, processConfig := range processConfigs {
		err = _startProcess(processConfig, data, comment, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
	- get persistence data entry for name
	  - if not present (assume not running), start process
	  - if present, check state
	    - state: running -> do nothing (already running)
		- state: stopped unexpected -> start process
*/
func _startProcess(processConfig *config.ProcessConfig, trackedData *persistence.Data, comment string, depLevel int) error {
	for _, dependencyName := range processConfig.DependsOn {
		if dependencyName != processConfig.Name {
			dc := CurrentContext.config.FindByName(dependencyName)
			if dc == nil {
				return fmt.Errorf("unable to find dependencyName '%s' for process '%s'", dependencyName, processConfig.Name)
			}
			err := _startProcess(dc, trackedData, comment, depLevel+1)
			if err != nil {
				return err
			}
		}
	}

	isAlreadyRunning, err := _isAlreadyRunning(processConfig.Name, trackedData)
	if err != nil {
		return err
	}
	startMessage := "Starting process"
	if depLevel > 0 {
		startMessage = "Starting dependency"
	}
	if depLevel == 0 || !isAlreadyRunning {
		_ = output.PrintMessageAndStatus(fmt.Sprintf("%s '%s'", startMessage, processConfig.Name), func() output.StatusReturn {
			dataEntry := trackedData.FindByName(processConfig.Name)
			if dataEntry != nil {
				pidExists, err := gopsutil.PidExists(dataEntry.Pid)
				if err != nil {
					return output.StatusReturn{Error: err}
				}
				if pidExists {
					return output.StatusReturn{OkMessage: "was already running"}
				}
			}
			_process := &process.Process{Config: processConfig}
			err := _process.Start()
			if err != nil {
				return output.StatusReturn{Error: err}
			}
			dataEntry, err = persistence.NewDataEntryFromProcess(_process)
			if err != nil {
				return output.StatusReturn{Error: err}
			}
			dataEntry.Comment = comment
			trackedData.AddOrUpdateEntry(dataEntry)
			return output.StatusReturn{Error: CurrentContext.persistenceWriter.Write(trackedData)}
		})
	}
	return nil
}

func _isAlreadyRunning(name string, trackedData *persistence.Data) (bool, error) {
	dataEntry := trackedData.FindByName(name)
	if dataEntry != nil {
		pidExists, err := gopsutil.PidExists(dataEntry.Pid)
		if err != nil {
			return false, err
		}
		return pidExists, nil
	}
	return false, nil
}
