package stop_strategy

import (
	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
	"syscall"
)

type StopStrategy interface {
	Stop(*config.ProcessConfig, *gopsutil.Process) error
}

func NewStopStrategyFromConfig(c *config.StopStrategyConfig) StopStrategy {
	if c == nil {
		return _getDefaultStopStrategy()
	}
	if c.Script != nil {
		return &ScriptStopStrategy{
			ScriptStopStrategyConfig: *c.Script,
		}
	} else if c.Signal != nil {
		return &SignalStopStrategy{
			SignalStopStrategyConfig: *c.Signal,
		}
	} else {
		return _getDefaultStopStrategy()
	}
}

func _getDefaultStopStrategy() StopStrategy {
	return &SignalStopStrategy{
		SignalStopStrategyConfig: config.SignalStopStrategyConfig{
			Signal: syscall.SIGTERM,
		},
	}
}
