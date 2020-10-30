package app

import (
	"fmt"
)

func InfoCommand(names []string, format string, filters Filters, columns []string) error {
	o := FormatMap[format]
	if o == nil {
		return fmt.Errorf("unknown format: %s", format)
	}
	o.SetWriter(CurrentContext.OutputWriter)
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, true)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}
	return o.Write(processes, columns)
}
