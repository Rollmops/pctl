package process

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	gopsutil "github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
)

func init() {
	s := &CommandlineEndsWithPidRetrieveStrategy{}
	PidRetrieveStrategies["command-ends-with"] = s
	PidRetrieveStrategies["cmd-ends-with"] = s
	PidRetrieveStrategies["cmdline-ends-with"] = s
}

type CommandlineEndsWithPidRetrieveStrategy struct{}

func (s *CommandlineEndsWithPidRetrieveStrategy) Retrieve(p *Process) (int32, error) {
	log.Tracef("Retrieving Pid from command end for %s", p.Config.Command)
	var pid int32

	if err := common.WaitUntilTrue(func() bool {
		var err error
		pid, err = _findPidForEndOfCommandline(p.Config.Command)
		if err != nil || pid == -1 {
			return false
		}
		return true
	}, CommandlinePidRetrieveStrategyWaitTime, CommandlinePidRetrieveStrategyAttempts); err != nil {
		return -1, err
	}

	return pid, nil
}

func _findPidForEndOfCommandline(command []string) (int32, error) {
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
		if len(command) > 0 && len(processCommand) > 0 &&
			strings.HasSuffix(strings.Join(processCommand, " "), strings.Join(command, " ")) {
			return _p.Pid, nil
		}
	}

	return -1, fmt.Errorf("unable to find process for command '%s'", command)

}
