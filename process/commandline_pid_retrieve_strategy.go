package process

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	gopsutil "github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

func init() {
	s := &CommandlinePidRetrieveStrategy{}
	PidRetrieveStrategies["command"] = s
	PidRetrieveStrategies["cmd"] = s
	PidRetrieveStrategies["cmdline"] = s
}

var (
	CommandlinePidRetrieveStrategyWaitTime      = 100 * time.Millisecond
	CommandlinePidRetrieveStrategyAttempts uint = 10
)

type CommandlinePidRetrieveStrategy struct{}

func (s *CommandlinePidRetrieveStrategy) Retrieve(p *Process) (int32, error) {
	log.Tracef("Retrieving Pid from command for %s", p.Config.Command)
	var pid int32

	if err := common.WaitUntilTrue(func() bool {
		var err error
		log.Tracef("Trying to find Pid for command %v", p.Config.Command)
		pid, err = _findPidForCommandline(p.Config.Command)
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
