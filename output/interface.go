package output

import (
	gopsutil "github.com/shirou/gopsutil/process"
)

type Output interface {
	Write([]*InfoEntry) error
}

type InfoEntry struct {
	Name           string
	ConfigCommand  string
	RunningCommand string
	// simple flag if process is running
	IsRunning bool
	// process was stopped, but not by pctl
	StoppedUnexpectedly bool
	// command of config differs from running command
	ConfigCommandChanged bool
	RunningInfo          *gopsutil.Process
}
