package app

import (
	"fmt"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/process"
	log "github.com/sirupsen/logrus"
)

func StopCommand(names []string, noWait bool, waitTime int) error {
	data, err := CurrentContext.persistenceReader.Read()
	if err != nil {
		return err
	}
	for _, name := range names {
		processConfig := CurrentContext.config.FindByName(name)
		if processConfig == nil {
			return fmt.Errorf("unable to find process '%s' in config", name)
		}
		dataEntry := data.FindByName(processConfig.Name)
		if dataEntry == nil {
			// TODO warn if we find a process with the same cmdline
			log.Infof("Process '%s' is not running", processConfig.Name)
			continue
		} else {
			p := process.Process{Config: processConfig}
			err = p.SynchronizeWithPid(dataEntry.Pid)
			if err != nil {
				return err
			}
			if !p.IsRunning() {
				log.Warnf("Expected '%s' as running ... no need to stop it", name)
			} else {
				log.Infof("Stopping process '%s'", processConfig.Name)

				err := output.PrintMessageAndStatus(fmt.Sprintf("Stopping process '%s'", processConfig.Name),
					func() error {
						return p.Stop()
					})
				if err != nil {
					return err
				}

				if !noWait {
					err = output.PrintMessageAndStatus(fmt.Sprintf("Waiting for stopped process '%s'", processConfig.Name),
						func() error {
							return p.WaitForStop(waitTime)
						},
					)
				}
				if err != nil {
					return err
				}
				data.RemoveByName(processConfig.Name)
			}
			err = CurrentContext.persistenceWriter.Write(data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
