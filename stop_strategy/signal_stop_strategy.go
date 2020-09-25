package stop_strategy

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
	"syscall"
)

var _signalNameMapping = map[string]syscall.Signal{
	"SIGTERM": syscall.SIGTERM,
	"SIGKILL": syscall.SIGKILL,
	"SIGINT":  syscall.SIGINT,
	"SIGQUIT": syscall.SIGQUIT,
	"SIGHUP":  syscall.SIGHUP,
}

type SignalStopStrategy struct {
	config.SignalStopStrategyConfig
}

func (s *SignalStopStrategy) Stop(_ string, p *gopsutil.Process) error {
	if s.SignalString != "" {
		signal, prs := _signalNameMapping[s.SignalString]
		if prs == false {
			return fmt.Errorf("invalid signal name: '%s'", s.SignalString)
		}
		return p.SendSignal(signal)
	}
	signal := s.Signal
	if signal == 0 {
		signal = syscall.SIGTERM
	}
	return p.SendSignal(signal)
}
