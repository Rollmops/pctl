package app

import (
	"time"
)

type ScriptProbeConfig struct {
	Script `yaml:",inline"`
}

type HttpProbeConfig struct {
}

type ProbeConfig struct {
	Timeout  string
	Interval string
	Script   *ScriptProbeConfig
}

func (c *ProbeConfig) Probe(process *Process) (bool, error) {
	if c.Script != nil {
		return c.Script.Probe(process)
	}

	return true, nil
}

func (c *ProbeConfig) GetTimeout() (time.Duration, error) {
	if c.Timeout == "" {
		return 5 * time.Second, nil
	}
	return time.ParseDuration(c.Timeout)
}

func (c *ProbeConfig) GetInterval() (time.Duration, error) {
	if c.Interval == "" {
		return 500 * time.Millisecond, nil
	}
	return time.ParseDuration(c.Interval)
}

type Probe interface {
	Probe(*Process) (bool, error)
}

func (s *ScriptProbeConfig) Probe(process *Process) (bool, error) {
	cmd, err := s.CreateCommand(process)
	if err != nil {
		return false, err
	}

	err = cmd.Run()
	return cmd.ProcessState.ExitCode() == 0, err
}

type NullProbe struct{}

func (p *NullProbe) Probe(_ *Process) (bool, error) {
	return true, nil
}
