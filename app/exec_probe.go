package app

import (
	"fmt"
	"time"
)

type ExecProbe struct {
	Exec `yaml:",inline"`
}

func (s *ExecProbe) Probe(process *Process, p *Probe) (bool, error) {
	timeoutDuration, err := p.GetTimeout()
	if err != nil {
		return false, err
	}
	cmd, err := s.CreateCommand(process)
	if err != nil {
		return false, err
	}

	err = cmd.Start()
	if err != nil {
		return false, err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	timeout := time.After(timeoutDuration)

	select {
	case <-timeout:
		_ = cmd.Process.Kill()
		return false, fmt.Errorf("exec command probe timed out")
	case err := <-done:
		return cmd.ProcessState.ExitCode() == 0, err
	}
}
