package app

import (
	"fmt"
	"github.com/Rollmops/pctl/output"
)

func ListCommand(format string) error {

	o := output.FormatMap[format]
	if o == nil {
		return fmt.Errorf("unknown output format: '%s'", format)
	}

	persistenceData, err := CurrentContext.persistenceReader.Read()
	if err != nil {
		return err
	}

	infoEntries, err := output.CreateInfoEntries(persistenceData, CurrentContext.config.Processes)
	if err != nil {
		return err
	}
	return o.Write(infoEntries)
}
