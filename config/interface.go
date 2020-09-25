package config

import (
	"syscall"
)

type Loader interface {
	Load(path string) (*Config, error)
}

type StopStrategyConfig struct {
	Script *ScriptStopStrategyConfig
	Signal *SignalStopStrategyConfig
}

type SignalStopStrategyConfig struct {
	Signal       syscall.Signal
	SignalString string
}

type ScriptStopStrategyConfig struct {
	Path          string
	Args          []string
	ForwardStdout bool `yaml:"forwardStdout"`
	ForwardStderr bool `yaml:"forwardStderr"`
}

type ProcessConfig struct {
	Name                    string
	Command                 []string            `yaml:"cmd"`
	PidRetrieveStrategyName string              `yaml:"pidStrategy"`
	StopStrategy            *StopStrategyConfig `yaml:"stopStrategy"`
}

type Config struct {
	Processes []*ProcessConfig
}
