package process

import "syscall"

type TerminateStrategy interface {
	Terminate(process *Process) error
}

type SignalTerminateStrategy struct {
	Signal syscall.Signal
}

func (s *SignalTerminateStrategy) Terminate(p *Process) error {
	info, err := p.Info()
	if err != nil {
		return err
	}
	return info.SendSignal(s.Signal)
}
