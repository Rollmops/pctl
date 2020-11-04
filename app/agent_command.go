package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

func AgentReloadCommand() error {
	InitializeAgentPidFile()
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
	InitializeAgentPidFile()
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

func AgentStartCommand(detach bool) error {
	InitializeAgentPidFile()
	if detach {
		err := SetAgentLogger()
		if err != nil {
			return err
		}
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
		logrus.Infof("Agent running on pid %d", pid)
	} else {
		return StartAgent()
	}
	return nil
}

func SetAgentLogger() error {
	if CurrentContext.Config.Agent.Logger != nil {
		if CurrentContext.Config.Agent.Logger.MaxBackups == 0 {
			CurrentContext.Config.Agent.Logger.MaxBackups = 5
		}
		logrus.SetOutput(CurrentContext.Config.Agent.Logger)
		if CurrentContext.Config.Agent.Logger.Level != "" {
			level, err := logrus.ParseLevel(CurrentContext.Config.Agent.Logger.Level)
			if err != nil {
				return err
			}
			logrus.SetLevel(level)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
	}
	return nil
}

func AgentStopCommand() error {
	InitializeAgentPidFile()
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
		logrus.Warningf("pctl agent is not running")
	}
	return nil
}
