package test

import (
	"bytes"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func init() {
	_ = SetConfigEnvPath("integration.yaml")
}

func SetConfigEnvPath(p ...string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pathAsList := []string{cwd, "..", "fixtures"}
	pathAsList = append(pathAsList, p...)
	configPath := path.Join(pathAsList...)
	return os.Setenv("PCTL_CONFIG_PATH", configPath)
}

func IsCommandRunning(command string) bool {
	processes, _ := gopsutil.Processes()
	for _, p := range processes {
		cmdline, _ := p.Cmdline()
		if cmdline == command {
			return true
		}
	}
	return false
}

func CaptureLogOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	logrus.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	logrus.SetOutput(os.Stderr)
	return buf.String()
}

func CaptureStdout(f func()) string {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	os.Stdout = w
	f()
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	return string(out)
}
