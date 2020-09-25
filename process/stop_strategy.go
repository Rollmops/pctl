package process

import "syscall"

type StopStrategy interface {
	Stop(process *Process) error
}

type SignalStopStrategy struct {
	Signal syscall.Signal
}

func (s *SignalStopStrategy) Stop(p *Process) error {
	info, err := p.Info()
	if err != nil {
		return err
	}
	return info.SendSignal(s.Signal)
}
