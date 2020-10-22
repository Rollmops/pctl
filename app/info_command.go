package app

import (
	"fmt"
	"github.com/Rollmops/pctl/output"
)

func InfoCommand(names []string, format string, filters []string) error {
	o := output.FormatMap[format]
	if o == nil {
		return fmt.Errorf("unknown format: %s", format)
	}
	o.SetWriter(CurrentContext.OutputWriter)

	processConfigs, err := CurrentContext.Config.CollectProcessConfigsByNameSpecifiers(names, filters, true)
	if err != nil {
		return err
	}

	infoEntries, err := output.CreateInfoEntries(processConfigs)
	if err != nil {
		return err
	}
	return o.Write(infoEntries)
}
