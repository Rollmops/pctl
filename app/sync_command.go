package app

import (
	"fmt"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/process"
	gopsutil "github.com/shirou/gopsutil/process"
	"strings"
)

var StrategyMapping = map[string]string{
	"exact":     "command",
	"ends-with": "command-ends-with",
}

func SyncCommand(names []string, strategy string) error {
	strategyName := StrategyMapping[strategy]
	if strategyName == "" {
		return fmt.Errorf("unknown sync strategy name: '%s'", strategy)
	}

	processConfigs := CurrentContext.config.CollectProcessConfigsByNameSpecifiers(names, true)
	if len(processConfigs) == 0 {
		return fmt.Errorf("no matching process config for name specifiers: %s", strings.Join(names, ", "))
	}

	// TODO change hacky implementation
	pidRetrieveStrategy := process.PidRetrieveStrategies[strategyName]
	process.CommandlinePidRetrieveStrategyAttempts = 1
	defer func() { process.CommandlinePidRetrieveStrategyAttempts = 10 }()

	data, err := CurrentContext.persistenceReader.Read()
	if err != nil {
		return nil
	}

	for _, processConfig := range processConfigs {
		_ = output.PrintMessageAndStatus(fmt.Sprintf("Syncing process '%s'", processConfig.Name),
			func() output.StatusReturn {
				dataEntry := data.FindByName(processConfig.Name)
				if dataEntry != nil {
					pidExists, err := gopsutil.PidExists(dataEntry.Pid)
					if err != nil {
						return output.StatusReturn{Error: err}
					}
					if pidExists {
						return output.StatusReturn{OkMessage: fmt.Sprintf("already tracked with PID: %d", dataEntry.Pid)}
					}
				}
				p := &process.Process{
					Config: processConfig,
				}
				pid, err := pidRetrieveStrategy.Retrieve(p)
				if err == nil && pid != -1 {
					dataEntry := &persistence.DataEntry{Pid: pid, Name: processConfig.Name, Command: processConfig.Command}
					data.AddOrUpdateEntry(dataEntry)
					err = CurrentContext.persistenceWriter.Write(data)
					if err != nil {
						return output.StatusReturn{Error: err}
					}
					return output.StatusReturn{OkMessage: fmt.Sprintf("command found with PID: %d", pid)}

				} else {
					return output.StatusReturn{Error: fmt.Errorf("no matching command found")}
				}
			})
	}
	return nil
}
