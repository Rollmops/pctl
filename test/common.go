package test

import gopsutil "github.com/shirou/gopsutil/process"

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
