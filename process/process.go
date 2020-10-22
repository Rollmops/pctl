package process

import (
	"encoding/json"
	"fmt"
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/stop_strategy"
	"os"
	"os/exec"
	"time"

	gopsutil "github.com/shirou/gopsutil/process"
)

type Process struct {
	Config        *config.ProcessConfig
	psutilProcess *gopsutil.Process
	cmd           *exec.Cmd
	environ       map[string]string
}

func (p *Process) SynchronizeWithPid(pid int32) error {
	_p, err := gopsutil.NewProcess(pid)
	if err != nil {
		return nil
	}
	p.psutilProcess = _p
	return nil
}

func (p *Process) IsRunning() bool {
	psutilProcess, err := p.GetPsutilProcess()
	if err != nil || psutilProcess == nil {
		return false
	}
	isRunning, err := psutilProcess.IsRunning()
	return isRunning && err == nil
}

func (p *Process) Pid() (int32, error) {
	runningEnvironInfo, err := FindRunningInfo(p.Config.Name)
	if err != nil {
		return -1, err
	}
	if runningEnvironInfo != nil {
		return runningEnvironInfo.Pid, nil
	}
	return -1, nil
}

func (p *Process) Start(comment string) error {
	name := p.Config.Command[0]

	var args []string
	if len(p.Config.Command) > 1 {
		args = p.Config.Command[1:]
	}

	runningInfoStr, err := _createRunningInfoJson(comment, p)
	if err != nil {
		return err
	}
	infoEnv := fmt.Sprintf("__PCTL_INFO__=%s", runningInfoStr)

	p.cmd = exec.Command(name, args...)
	p.cmd.Env = os.Environ()
	p.cmd.Env = append(p.cmd.Env, infoEnv)

	return p.cmd.Start()
}

func _createRunningInfoJson(comment string, p *Process) (string, error) {
	runningEnvironInfo := RunningEnvironInfo{
		Config:  *p.Config,
		Comment: comment,
	}
	infoStr, err := json.Marshal(runningEnvironInfo)
	if err != nil {
		return "", err
	}
	return string(infoStr), nil
}

func (p *Process) WaitForStarted(maxWaitTime time.Duration, intervalDuration time.Duration) error {
	attempts := maxWaitTime / intervalDuration
	err := common.WaitUntilTrue(func() bool {
		return p.IsRunning()
	}, intervalDuration, uint(attempts))
	if err != nil {
		return fmt.Errorf("unable to start process '%s'", p.Config.Name)
	}
	return nil
}

func (p *Process) Stop() error {
	_process, err := p.GetPsutilProcess()
	if err != nil {
		return err
	}
	stopStrategy := stop_strategy.NewStopStrategyFromConfig(p.Config.StopStrategy)
	err = stopStrategy.Stop(p.Config, _process)
	if err != nil {
		return err
	}
	return nil
}

func (p *Process) WaitForStop(maxWaitTime time.Duration, intervalDuration time.Duration) error {
	attempts := maxWaitTime / intervalDuration
	err := common.WaitUntilTrue(func() bool {
		return !p.IsRunning()
	}, intervalDuration, uint(attempts))
	if err != nil {
		pid, _ := p.Pid()
		return fmt.Errorf("unable to stop process '%s' on PID %d", p.Config.Name, pid)
	}
	return nil
}

func (p *Process) Kill() error {
	return p.cmd.Process.Kill()
}

func (p *Process) GetPsutilProcess() (*gopsutil.Process, error) {
	if p.psutilProcess != nil {
		return p.psutilProcess, nil
	}
	pid, err := p.Pid()
	if err != nil {
		return nil, err
	}
	if pid == -1 {
		return nil, nil
	}
	psutilProcess, err := gopsutil.NewProcess(pid)
	p.psutilProcess = psutilProcess
	return psutilProcess, err
}
