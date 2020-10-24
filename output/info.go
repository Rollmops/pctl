package output

import (
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
	gopsutil "github.com/shirou/gopsutil/process"
)

type Info struct {
	Name                 string
	Comment              string
	ConfigCommand        []string
	RunningCommand       []string
	IsRunning            bool
	ConfigCommandChanged bool
	RunningInfo          *gopsutil.Process
}

func CreateInfos(processConfigs []*config.ProcessConfig) ([]*Info, error) {
	var infos []*Info
	for _, processConfig := range processConfigs {
		infoEntry := &Info{
			Name:           processConfig.Name,
			ConfigCommand:  processConfig.Command,
			RunningCommand: processConfig.Command,
		}

		runningProcessConfig, err := process.FindRunningInfo(processConfig.Name)
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

		infos = append(infos, infoEntry)
	}
	return infos, nil
}
