package output

import (
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
)

func CreateInfoEntries(processConfigs []*config.ProcessConfig) ([]*InfoEntry, error) {
	var infoEntries []*InfoEntry
	for _, processConfig := range processConfigs {
		infoEntry := &InfoEntry{
			Name:           processConfig.Name,
			ConfigCommand:  processConfig.Command,
			RunningCommand: processConfig.Command,
		}

		runningProcessConfig, err := process.FindRunningEnvironInfoFromName(processConfig.Name)
		if err != nil {
			return nil, err
		}
		if runningProcessConfig == nil {
			infoEntry.IsRunning = false
			infoEntry.ConfigCommand = processConfig.Command

		} else {
			infoEntry.IsRunning = true
			infoEntry.ConfigCommand = runningProcessConfig.Config.Command
			infoEntry.ConfigCommandChanged = !common.CompareStringSlices(infoEntry.ConfigCommand, processConfig.Command)
			p := process.Process{Config: processConfig}
			err = p.SynchronizeWithPid(runningProcessConfig.Pid)
			if err != nil {
				return nil, err
			}
			infoEntry.RunningInfo, err = p.GetPsutilProcess()
			if err != nil {
				return nil, err
			}
			infoEntry.RunningCommand, err = infoEntry.RunningInfo.CmdlineSlice()
			if err != nil {
				return nil, err
			}
		}

		infoEntries = append(infoEntries, infoEntry)
	}
	return infoEntries, nil
}
