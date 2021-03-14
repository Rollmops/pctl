package app

import (
	"syscall"
	"time"
)

type StopConfig struct {
	Exec    *Exec                     `yaml:"exec"`
	Signal  *SignalStopStrategyConfig `yaml:"signal"`
	Timeout string                    `yaml:"timeout"`
	Period  string                    `yaml:"period"`
}

type StopStrategy interface {
	Stop(process *Process) error
}

func (c *StopConfig) GetTimeout() (time.Duration, error) {
	if c.Timeout == "" {
		return 5 * time.Second, nil
	}
	return time.ParseDuration(c.Timeout)
}

func (c *StopConfig) GetInterval() (time.Duration, error) {
	if c.Period == "" {
		return 10 * time.Millisecond, nil
	}
	return time.ParseDuration(c.Period)
}

func NewStopStrategyFromConfig(c *StopConfig) StopStrategy {
	if c == nil {
		return _getDefaultStopStrategy()
	}
	if c.Exec != nil {
		return &ScriptStopStrategy{
			Exec: *c.Exec,
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
