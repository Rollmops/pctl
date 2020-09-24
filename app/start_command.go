package app

import (
	"fmt"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	log "github.com/sirupsen/logrus"
)

func StartCommand(names []string) error {
	for _, name := range names {
		_config := CurrentContext.config.FindByName(name)
		if _config == nil {
			return fmt.Errorf("unable to find process '%s'", name)
		}
		_process := process.NewProcess(*_config)
		log.Infof("Starting process '%s'", name)
		err := _process.Start()
		if err != nil {
			return err
		}
		pid, err := _process.Pid()
		if err != nil {
			return err
		}
		info, err := _process.Info()
		if err != nil {
			return err
		}
		cmdline, err := info.Cmdline()
		if err != nil {
			return err
		}
		data, err := CurrentContext.persistenceReader.Read()
		data = append(data, persistence.Data{
			Name: name,
			Pid:  pid,
			Cmd:  cmdline,
		})

		err = CurrentContext.persistenceWriter.Write(data)
		if err != nil {
			return err
		}

	}
	return nil
}
