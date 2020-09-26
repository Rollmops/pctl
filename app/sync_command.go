package app

import (
	"fmt"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	gopsutil "github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

var StrategyMapping = map[string]string{
	"exact":     "command",
	"ends-with": "command_ends_with",
}

func SyncCommand(names []string, all bool, strategy string) error {
	if all {
		log.Debug("Synchronizing all processes")
		names = CurrentContext.config.GetAllProcessNames()
	}

	strategyName := StrategyMapping[strategy]
	if strategyName == "" {
		return fmt.Errorf("unknown sync strategy name: '%s'", strategy)
	}

	pidRetrieveStrategy := process.PidRetrieveStrategies[strategyName]
	process.CommandlinePidRetrieveStrategyAttempts = 1
	defer func() { process.CommandlinePidRetrieveStrategyAttempts = 10 }()

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
		if dataEntry != nil {
			pidExists, err := gopsutil.PidExists(dataEntry.Pid)
			if err != nil {
				return nil
			}
			if pidExists {
				log.Warningf("Process '%s' is already tracked and running", processConfig.Name)
				continue
			}
		}
		p := &process.Process{
			Config: processConfig,
		}
		pid, err := pidRetrieveStrategy.Retrieve(p)
		if err == nil && pid != -1 {
			fmt.Printf("Found matching command for process '%s' with PID %d\n", processConfig.Name, pid)
			dataEntry := &persistence.DataEntry{Pid: pid, Name: processConfig.Name, Command: processConfig.Command}
			data.AddOrUpdateEntry(dataEntry)
			err = CurrentContext.persistenceWriter.Write(data)
			if err != nil {
				return err
			}
		} else {
			log.Warningf("No matching command found for process '%s' and strategy '%s'\n", processConfig.Name, strategy)
		}
	}
	return nil
}
