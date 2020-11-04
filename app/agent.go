package app

import (
	"github.com/facebookgo/pidfile"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

type AgentLoggerConfig struct {
	lumberjack.Logger `yaml:",inline"`
	Level             string `yaml:"level"`
}

type AgentConfig struct {
	PidFile string             `yaml:"pidFile"`
	Logger  *AgentLoggerConfig `yaml:"logger"`
}

func FindAgentProcess() (*gopsutil.Process, error) {
	pid, err := pidfile.Read()
	if err != nil {
		return nil, nil
	}
	p, _ := gopsutil.NewProcess(int32(pid))
	isRunning, _ := p.IsRunning()
	if !isRunning {
		return nil, nil
	}
	return p, nil
}

func CheckRunningAgentProcess() error {
	agentProcess, err := FindAgentProcess()
	if err != nil {
		return err
	}
	if agentProcess != nil {
		logrus.Fatalf("Agent process already running on PID %d", agentProcess.Pid)
	}
	return nil
}

func StartAgent() error {
	err := CheckRunningAgentProcess()
	if err != nil {
		return err
	}
	err = pidfile.Write()
	if err != nil {
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGUSR1)
	go StartAgentMain()
	logrus.Infof("Starting %d watchers", len(CurrentContext.Config.ProcessConfigs))
	for _, processConfig := range CurrentContext.Config.ProcessConfigs {
		watcher := AgentWatcher{processConfig: processConfig}
		go watcher.Start()
	}

	for {
		s := <-c
		if s == syscall.SIGTERM {
			logrus.Infof("Received TERM signal")
			os.Exit(0)
		} else if s == syscall.SIGUSR1 {
			logrus.Infof("Reloading config")
			err := CurrentContext.InitializeRunningProcessInfo()
			if err != nil {
				logrus.Fatalf(err.Error())
			}
		}
	}
}

func StartAgentMain() {
	logrus.Infof("Starting agent main")
	for {
		time.Sleep(5 * time.Second)
		logrus.Debugf("Refreshing process info")
		err := CurrentContext.Cache.Refresh()
		if err != nil {
			logrus.Error(err.Error())
		}
	}
}

func InitializeAgentPidFile() {
	if CurrentContext.Config.Agent.PidFile != "" {
		pidfile.SetPidfilePath(CurrentContext.Config.Agent.PidFile)
	} else {
		pidfile.SetPidfilePath(path.Join("/", "tmp", "pctl-agent.pid"))
	}
}
