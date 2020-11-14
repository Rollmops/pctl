package app

import (
	"fmt"
	"github.com/facebookgo/pidfile"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/ioutil"
	"net"
	"net/http"
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

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

func hello(w http.ResponseWriter, req *http.Request) {
	fw := flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}
	bodyReader := req.Body
	defer bodyReader.Close()
	var buffer []byte

	buffer, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("%s\n", string(buffer))
	fw.Write([]byte("wuff\n"))
	time.Sleep(2 * time.Second)
	fw.Write([]byte("waff\n"))
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
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGUSR1)

	go StartAgentMain()
	logrus.Infof("Starting %d watchers", len(CurrentContext.Config.ProcessConfigs))
	for _, processConfig := range CurrentContext.Config.ProcessConfigs {
		watcher := AgentWatcher{processConfig: processConfig}
		go watcher.Start()
	}

	go func() {
		err = startServer()
		if err != nil {
			panic(err)
		}
	}()

	for {
		s := <-c
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			logrus.Infof("Received TERM signal")
			os.Remove("/tmp/pctl-agent.sock")
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

func startServer() error {
	http.HandleFunc("/pctl", hello)
	unixListener, err := net.Listen("unix", "/tmp/pctl-agent.sock")
	if err != nil {
		return err
	}
	err = http.Serve(unixListener, nil)
	if err != nil {
		return err
	}
	return nil
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
