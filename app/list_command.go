package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
)

func ListCommand(names []string, format string) error {

	o := output.FormatMap[format]
	if o == nil {
		return fmt.Errorf("unknown output format: '%s'", format)
	}

	persistenceData, err := CurrentContext.persistenceReader.Read()
	if err != nil {
		return err
	}

	var filteredProcessConfigs []*config.ProcessConfig
	if len(names) > 0 {
		for _, name := range names {
			c := CurrentContext.config.FindByName(name)
			if c == nil {
				return fmt.Errorf("unble to find process '%s'", name)
			}
			filteredProcessConfigs = append(filteredProcessConfigs, c)
		}
	} else {
		filteredProcessConfigs = CurrentContext.config.Processes
	}
	infoEntries, err := output.CreateInfoEntries(persistenceData, CurrentContext.config.Processes)
	if err != nil {
		return err
	}
	return o.Write(infoEntries)
}
