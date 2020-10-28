package app

import (
	"fmt"
	"github.com/facebookgo/pidfile"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func AgentReloadCommand() error {
	agentProcess, err := FindAgentProcess()
	if err != nil {
		return err
	}
	if agentProcess == nil {
		logrus.Fatal("Agent is not running")
	}
	logrus.Infof("Sending SIGUSR1 to %d", agentProcess.Pid)
	return agentProcess.SendSignal(syscall.SIGUSR1)
}

func AgentStatusCommand(deriveExitCode bool) error {
	agentProcess, err := FindAgentProcess()
	if err != nil {
		return err
	}
	if agentProcess == nil {
		fmt.Println("The agent is not running")
		if deriveExitCode {
			os.Exit(1)
		}
	} else {
		fmt.Printf("The agent is running with PID %d\n", agentProcess.Pid)
	}
	return nil
}

func AgentStartCommand(logPath string, detach bool) error {
	logrus.SetLevel(logrus.InfoLevel)
	if logPath != "" {
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}
	_, isAgent := os.LookupEnv("__PCTL_AGENT__")
	if !isAgent && detach {
		env := []string{"__PCTL_AGENT__=1"}
		pctlPath, err := exec.LookPath(os.Args[0])
		if err != nil {
			return err
		}
		pid, err := syscall.ForkExec(pctlPath, os.Args,
			&syscall.ProcAttr{
				Env: append(os.Environ(), env...),
				Sys: &syscall.SysProcAttr{
					Setsid: true,
				},
			})
		if err != nil {
			return err
		}
		logrus.Infof("Watcher running on pid %d", pid)
	} else {
		return startAgent()
	}
	return nil
}

func startAgent() error {
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
	logrus.Infof("Starting %d watchers", len(CurrentContext.Config.ProcessConfigs))
	for _, processConfig := range CurrentContext.Config.ProcessConfigs {
		watcher := AgentWatcher{processConfig: processConfig}
		go watcher.Start()
	}

	for {
		s := <-c
		if s == syscall.SIGTERM {
			logrus.Info("Received TERM signal")
			os.Exit(0)
		} else if s == syscall.SIGUSR1 {
			logrus.Info("Reloading config")
			err := CurrentContext.Initialize()
			if err != nil {
				logrus.Fatal(err)
			}
		}
	}
}

func AgentStopCommand() error {
	p, err := FindAgentProcess()
	if err != nil {
		return err
	}
	if p != nil {
		err = p.SendSignal(syscall.SIGTERM)
		if err != nil {
			return err
		}
	} else {
		logrus.Warning("pctl agent is not running")
	}
	return nil
}
