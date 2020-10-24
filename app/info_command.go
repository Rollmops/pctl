package app

import (
	"fmt"
)

func InfoCommand(names []string, format string, filters []string) error {
	o := FormatMap[format]
	if o == nil {
		return fmt.Errorf("unknown format: %s", format)
	}
	o.SetWriter(CurrentContext.OutputWriter)

	processes, err := CurrentContext.Config.CollectSyncedProcessesByNameSpecifiers(names, filters, true)
	if err != nil {
		return err
	}

	return o.Write(processes)
}
