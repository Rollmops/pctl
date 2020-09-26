package output

import (
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	"github.com/fatih/color"
)

var FormatMap = map[string]Output{}

var Green = color.New(color.FgGreen).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()

func CreateInfoEntries(persistenceData *persistence.Data, processConfigs []*config.ProcessConfig) ([]*InfoEntry, error) {
	var infoEntries []*InfoEntry
	for _, processConfig := range processConfigs {
		infoEntry := &InfoEntry{
			Name:           processConfig.Name,
			ConfigCommand:  processConfig.Command,
			RunningCommand: processConfig.Command,
		}

		if e := persistenceData.FindByName(processConfig.Name); e != nil {
			infoEntry.RunningCommand = e.Command
			p := process.Process{Config: processConfig}
			err := p.SynchronizeWithPid(e.Pid)
			if err != nil {
				return nil, err
			}
			if p.IsRunning() {
				infoEntry.RunningInfo, err = p.Info()
				if err != nil {
					return nil, err
				}
			}
			infoEntry.IsRunning = p.IsRunning()
			infoEntry.StoppedUnexpectedly = !infoEntry.IsRunning
			infoEntry.ConfigCommandChanged = !common.CompareStringSlices(e.Command, processConfig.Command)
		}
		infoEntries = append(infoEntries, infoEntry)
	}
	return infoEntries, nil
}
