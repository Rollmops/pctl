package app

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/process"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func StopCommand(names []string, noWait bool) error {
	processConfigs := CurrentContext.Config.CollectProcessConfigsByNameSpecifiers(names, false)
	if len(processConfigs) == 0 {
		return fmt.Errorf("no matching process Config for name specifiers: %s", strings.Join(names, ", "))
	}
	for _, processConfig := range processConfigs {
		for _, dependentOf := range CurrentContext.Config.GetAllDependentOf(processConfig.Name) {
			if dependentOf.Name != processConfig.Name {
				_ = output.PrintMessageAndStatus(fmt.Sprintf("Stopping dependency process '%s' of '%s", dependentOf.Name, processConfig.Name),
					func() output.StatusReturn {
						return _stopProcess(dependentOf, noWait)
					},
				)
			}
		}
		_ = output.PrintMessageAndStatus(fmt.Sprintf("Stopping process '%s", processConfig.Name),
			func() output.StatusReturn {
				return _stopProcess(processConfig, noWait)
			},
		)
	}
	return nil
}

func _stopProcess(processConfig *config.ProcessConfig, noWait bool) output.StatusReturn {
	runningEnvironInfo, err := process.FindRunningEnvironInfoFromName(processConfig.Name)
	if err != nil {
		return output.StatusReturn{Error: err}
	}

	if runningEnvironInfo == nil {
		// TODO warn if we find a process with the same cmdline
		return output.StatusReturn{OkMessage: "was not running"}
	} else {
		p := process.Process{Config: processConfig}
		err := p.SynchronizeWithPid(runningEnvironInfo.Pid)
		if err != nil {
			return output.StatusReturn{Error: err}
		}
		if !p.IsRunning() {
			return output.StatusReturn{WarningMessage: "tracked as running but stopped unexpectedly"}
		} else {
			err = p.Stop()
			if noWait {
				return output.StatusReturn{Error: err}
			}
			maxWaitTime, intervalTime, err := _getMaxWaitTimeAndIntervalDuration(processConfig)
			if err != nil {
				return output.StatusReturn{Error: err}
			}
			return output.StatusReturn{Error: p.WaitForStop(maxWaitTime, intervalTime)}
		}
	}
}

func _getMaxWaitTimeAndIntervalDuration(p *config.ProcessConfig) (time.Duration, time.Duration, error) {
	maxWaitTime := 5 * time.Second
	intervalTime := 100 * time.Millisecond
	if p.StopStrategy != nil {
		if p.StopStrategy.MaxWaitTime != "" {
			maxWaitTime, err := time.ParseDuration(p.StopStrategy.MaxWaitTime)
			if err != nil {
				return maxWaitTime, intervalTime, err
			}
		}
		if p.StopStrategy.IntervalTime != "" {
			intervalTime, err := time.ParseDuration(p.StopStrategy.IntervalTime)
			if err != nil {
				return maxWaitTime, intervalTime, err
			}
		}
	}
	log.Tracef("max wait time: %d, interval time: %d", maxWaitTime, intervalTime)
	return maxWaitTime, intervalTime, nil
}
