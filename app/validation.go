package app

import (
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
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
				logrus.Errorf(
					"Found tracked running process '%s' with pid %d that could not be found in config",
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
