package process

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	gopsutil "github.com/shirou/gopsutil/process"
	"sort"
	"time"
)

func init() {
	s := &CommandlinePidRetrieveStrategy{}
	PidRetrieveStrategies["command"] = s
}

var (
	CommandlinePidRetrieveStrategyWaitTime      = 10 * time.Millisecond
	CommandlinePidRetrieveStrategyAttempts uint = 100
)

type CommandlinePidRetrieveStrategy struct{}

func (s *CommandlinePidRetrieveStrategy) Retrieve(p *Process) (int32, error) {
	var pid int32

	if err := common.WaitUntilTrue(func() bool {
		var err error
		pid, err = _findPidForCommandline(p.cmd.Args)
		if err != nil || pid == -1 {
			return false
		}
		return true
	}, CommandlinePidRetrieveStrategyWaitTime, CommandlinePidRetrieveStrategyAttempts); err != nil {
		return -1, err
	}

	return pid, nil
}

func _findPidForCommandline(command []string) (int32, error) {
	processes, err := gopsutil.Processes()
	if err != nil {
		return -1, err
	}
	sort.SliceStable(processes, func(i, j int) bool {
		firstCreateTime, _ := processes[i].CreateTime()
		secondCreateTime, _ := processes[j].CreateTime()
		return firstCreateTime < secondCreateTime
	})

	for _, _p := range processes {
		processCommand, err := _p.CmdlineSlice()
		if err != nil {
			return -1, err
		}
		if common.CompareStringSlices(processCommand, command) {
			return _p.Pid, nil
		}
	}

	return -1, fmt.Errorf("unable to find process for command '%s'", command)

}
