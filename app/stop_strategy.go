package app

import (
	"syscall"
)

type StopStrategy interface {
	Stop(process *Process) error
}

func NewStopStrategyFromConfig(c *StopStrategyConfig) StopStrategy {
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
		SignalStopStrategyConfig: SignalStopStrategyConfig{
			Signal: syscall.SIGTERM,
		},
	}
}
