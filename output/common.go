package output

import (
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	"strings"
)

var FormatMap = map[string]Output{}

func CreateInfoEntries(persistenceData *persistence.Data, processConfigs []*config.ProcessConfig) ([]*InfoEntry, error) {
	var infoEntries []*InfoEntry
	for _, processConfig := range processConfigs {
		configCommand := strings.Join(processConfig.Cmd, " ")
		infoEntry := &InfoEntry{
			Name:           processConfig.Name,
			ConfigCommand:  configCommand,
			RunningCommand: configCommand,
		}

		if e := persistenceData.FindByName(processConfig.Name); e != nil {
			infoEntry.RunningCommand = e.Cmd
			p := process.NewProcess(processConfig)
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
			infoEntry.ConfigCommandChanged = e.Cmd != configCommand
		}
		infoEntries = append(infoEntries, infoEntry)
	}
	return infoEntries, nil
}
