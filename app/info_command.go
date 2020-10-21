package app

import (
	"fmt"
	"github.com/Rollmops/pctl/output"
)

func InfoCommand(names []string, format string) error {
	o := output.FormatMap[format]
	if o == nil {
		return fmt.Errorf("unknown format: %s", format)
	}
	o.SetWriter(CurrentContext.OutputWriter)

	processConfigs := CurrentContext.Config.CollectProcessConfigsByNameSpecifiers(names, true)

	infoEntries, err := output.CreateInfoEntries(processConfigs)
	if err != nil {
		return err
	}
	return o.Write(infoEntries)
}
