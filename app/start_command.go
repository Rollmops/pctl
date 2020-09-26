package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	gopsutil "github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

func StartCommand(names []string, all bool, comment string) error {

	if all {
		log.Debug("Starting all processes")
		names = CurrentContext.config.GetAllProcessNames()
	}
	log.Tracef("Starting processes: %v", names)
	data, err := CurrentContext.persistenceReader.Read()
	if err != nil {
		return nil
	}
	for _, name := range names {
		processConfig := CurrentContext.config.FindByName(name)
		if processConfig == nil {
			return fmt.Errorf("unable to find process '%s' in config", name)
		}
		err = _startProcess(processConfig, data, comment)
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
func _startProcess(processConfig *config.ProcessConfig, data *persistence.Data, comment string) error {
	log.Debugf("Starting process '%s'", processConfig.Name)
	dataEntry := data.FindByName(processConfig.Name)
	var startNeeded bool
	if dataEntry != nil {
		log.Tracef("persistence data entry was found")
		pidExists, err := gopsutil.PidExists(dataEntry.Pid)
		if err != nil {
			return err
		}
		if !pidExists {
			log.Tracef("Pid %d is running for persistence data entry", dataEntry.Pid)
			log.Warnf("Expected not running process '%s' as running ... starting it", processConfig.Command)
			startNeeded = true
		} else {
			startNeeded = false
			fmt.Printf("Process '%s' is already running\n", processConfig.Name)
		}
	} else {
		startNeeded = true
	}
	log.Tracef("Start needed: %v", startNeeded)
	if len(processConfig.DependsOn) > 0 {
		fmt.Println("Starting dependencies")
	}
	for _, dependencyName := range processConfig.DependsOn {
		dc := CurrentContext.config.FindByName(dependencyName)
		if dc == nil {
			return fmt.Errorf("unable to find dependencyName '%s' for process '%s'", dependencyName, processConfig.Name)
		}
		err := _startProcess(dc, data, comment)
		if err != nil {
			return err
		}
	}
	if startNeeded {
		_process := &process.Process{Config: processConfig}

		err := output.PrintMessageAndStatus(fmt.Sprintf("Starting process '%s'", processConfig.Name),
			func() error {
				err := _process.Start()
				if err != nil {
					return err
				}
				dataEntry, err = persistence.NewDataEntryFromProcess(_process)
				if err != nil {
					return err
				}
				dataEntry.Comment = comment
				data.AddOrUpdateEntry(dataEntry)
				return CurrentContext.persistenceWriter.Write(data)
			})
		if err != nil {
			return err
		}
	}
	return nil
}
