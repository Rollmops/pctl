package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
)

func InfoCommand(names []string, format string) error {
	o := output.FormatMap[format]
	o.SetWriter(CurrentContext.OutputWriter)
	if o == nil {
		return fmt.Errorf("unknown Output format: '%s'", format)
	}

	persistenceData, err := CurrentContext.PersistenceReader.Read()
	if err != nil {
		return err
	}

	var filteredProcessConfigs []*config.ProcessConfig
	if len(names) > 0 {
		for _, name := range names {
			c := CurrentContext.Config.FindByName(name)
			if c == nil {
				return fmt.Errorf("unble to find process '%s'", name)
			}
			filteredProcessConfigs = append(filteredProcessConfigs, c)
		}
	} else {
		filteredProcessConfigs = CurrentContext.Config.Processes
	}
	infoEntries, err := output.CreateInfoEntries(persistenceData, CurrentContext.Config.Processes)
	if err != nil {
		return err
	}
	return o.Write(infoEntries)
}
