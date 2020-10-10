package app

import (
	"fmt"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"github.com/yourbasic/graph"
)

func ValidatePersistenceConfigDiscrepancy() error {
	logrus.Debug("Checking for config - persistence discrepancies")
	data, err := CurrentContext.PersistenceReader.Read()
	if err != nil {
		return err
	}
	for _, dataEntry := range data.Entries {
		if p := CurrentContext.Config.FindByName(dataEntry.Name); p == nil {
			isRunning, err := gopsutil.PidExists(dataEntry.Pid)
			if err != nil {
				return err
			}
			if isRunning {
				return fmt.Errorf("found tracked running process '%s' with pid %d that could not be found in config",
					dataEntry.Name, dataEntry.Pid)
			} else {
				logrus.Warningf("Found tracked process '%s' that is not running and not found in config - removing it",
					dataEntry.Name)
				data.RemoveByName(dataEntry.Name)
				err = CurrentContext.PersistenceWriter.Write(data)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func ValidateAcyclicDependencies() error {
	mapping := make(map[string]int)
	for i, _config := range CurrentContext.Config.Processes {
		mapping[_config.Name] = i
	}

	gm := graph.New(len(CurrentContext.Config.Processes))
	for _, _config := range CurrentContext.Config.Processes {
		for _, n := range _config.DependsOn {
			gm.Add(mapping[_config.Name], mapping[n])
		}
	}
	if !graph.Acyclic(gm) {
		return fmt.Errorf("process dependency configuration is not acyclic")
	}
	return nil
}
