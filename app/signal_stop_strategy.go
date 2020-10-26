package app

import (
	"fmt"
	"syscall"
)

type SignalStopStrategyConfig struct {
	Signal       syscall.Signal
	SignalString string
}

var _signalNameMapping = map[string]syscall.Signal{
	"SIGTERM": syscall.SIGTERM,
	"SIGKILL": syscall.SIGKILL,
	"SIGINT":  syscall.SIGINT,
	"SIGQUIT": syscall.SIGQUIT,
	"SIGHUP":  syscall.SIGHUP,
}

type SignalStopStrategy struct {
	SignalStopStrategyConfig
}

func (s *SignalStopStrategy) Stop(p *Process) error {
	if s.SignalString != "" {
		signal, prs := _signalNameMapping[s.SignalString]
		if prs == false {
			return fmt.Errorf("invalid signal name: '%s'", s.SignalString)
		}
		return p.Info.GoPsutilProcess.SendSignal(signal)
	}
	signal := s.Signal
	if signal == 0 {
		signal = syscall.SIGTERM
	}
	return p.Info.GoPsutilProcess.SendSignal(signal)
}
