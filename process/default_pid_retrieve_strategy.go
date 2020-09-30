package process

import (
	"fmt"
)

func init() {
	PidRetrieveStrategies[""] = &DefaultPidRetrieveStrategy{}
	PidRetrieveStrategies["default"] = &DefaultPidRetrieveStrategy{}
}

type DefaultPidRetrieveStrategy struct{}

func (s *DefaultPidRetrieveStrategy) Retrieve(p *Process) (int32, error) {
	if p.cmd == nil || p.cmd.Process == nil {
		return -1, fmt.Errorf("command or process object is nil")
	}

	return int32(p.cmd.Process.Pid), nil
}
