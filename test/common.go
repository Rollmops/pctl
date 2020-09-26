package test

import (
	"bytes"
	"github.com/Rollmops/pctl/app"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func init() {
	cwd, _ := os.Getwd()
	configPath := path.Join(cwd, "..", "fixtures", "integration.yaml")
	_ = os.Setenv("PCTL_CONFIG_PATH", configPath)
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

func StartAppAndGetStdout(args []string) (string, error) {
	r, w, _ := os.Pipe()

	pctlApp := app.CreateCliApp(w)
	err := pctlApp.Run(args)
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
