package stop_strategy

import (
	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
	"syscall"
)

type StopStrategy interface {
	Stop(string, *gopsutil.Process) error
}

func NewStopStrategyFromConfig(c *config.StopStrategyConfig) StopStrategy {
	if c == nil {
		return &SignalStopStrategy{SignalStopStrategyConfig: config.SignalStopStrategyConfig{Signal: syscall.SIGTERM}}
	}
	if c.Script != nil {
		return &ScriptStopStrategy{ScriptStopStrategyConfig: *c.Script}
	} else {
		return &SignalStopStrategy{SignalStopStrategyConfig: *c.Signal}
	}
}
