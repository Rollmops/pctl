package app

import (
	"fmt"
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
		dataEntry := data.FindByName(name)
		if dataEntry == nil {
			// TODO warn if we find a process with the same cmdline
			_process, err := _startProcessByName(name)
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
				_process, err := _startProcessByName(name)
				if err != nil {
					return err
				}
				dataEntry, err = persistence.NewDataEntryFromProcess(_process)
				if err != nil {
					return err
				}
			}
		}
		data.AddOrUpdateEntry(dataEntry)
	}

	return CurrentContext.persistenceWriter.Write(data)
}

func _startProcessByName(name string) (*process.Process, error) {
	_config := CurrentContext.config.FindByName(name)
	if _config == nil {
		return nil, fmt.Errorf("unable to find process '%s'", name)
	}
	log.Infof("Starting process '%s'", name)
	_process := process.NewProcess(*_config)
	err := _process.Start()
	if err != nil {
		return nil, err
	}
	return _process, nil
}
