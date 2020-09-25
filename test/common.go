package test

import (
	gopsutil "github.com/shirou/gopsutil/process"
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
