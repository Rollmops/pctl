package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	gopsutil "github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

func StartCommand(names []string) error {
	/*
		- get persistence data entry for name
		  - if not present (assume not running), start process
		  - if present, check state
		    - state: running -> do nothing (already running)
			- state: stopped unexpected -> start process
	*/
	data, err := CurrentContext.persistenceReader.Read()
	if err != nil {
		return nil
	}
	for _, name := range names {
		processConfig := CurrentContext.config.FindByName(name)
		if processConfig == nil {
			return fmt.Errorf("unable to find process '%s' in config", name)
		}
		dataEntry := data.FindByName(processConfig.Name)
		if dataEntry == nil {
			// TODO warn if we find a process with the same cmdline
			_process, err := _startProcess(processConfig)
			if err != nil {
				return err
			}
			dataEntry, err = persistence.NewDataEntryFromProcess(_process)
			if err != nil {
				return err
			}
		} else {
			pidExists, err := gopsutil.PidExists(dataEntry.Pid)
			if err != nil {
				return err
			}
			if !pidExists {
				log.Warnf("Expected '%s' as running ... starting it", name)
				_process, err := _startProcess(processConfig)
				if err != nil {
					return err
				}
				dataEntry, err = persistence.NewDataEntryFromProcess(_process)
				if err != nil {
					return err
				}
			} else {
				log.Infof("Process '%s' is already running", name)
			}
		}
		data.AddOrUpdateEntry(dataEntry)
	}

	return CurrentContext.persistenceWriter.Write(data)
}

func _startProcess(processConfig *config.ProcessConfig) (*process.Process, error) {
	log.Infof("Starting process '%s'", processConfig.Name)
	_process := &process.Process{Config: processConfig}
	err := _process.Start()
	if err != nil {
		return nil, err
	}
	return _process, nil
}
