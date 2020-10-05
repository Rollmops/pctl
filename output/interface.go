package output

import (
	gopsutil "github.com/shirou/gopsutil/process"
	"os"
)

type Output interface {
	Write([]*InfoEntry) error
	SetWriter(file *os.File)
}

type InfoEntry struct {
	Name           string
	Comment        string
	ConfigCommand  []string
	RunningCommand []string
	MarkFlag       int
	// simple flag if process is running
	IsRunning bool
	// process was stopped, but not by pctl
	StoppedUnexpectedly bool
	// command of config differs from running command
	ConfigCommandChanged bool
	RunningInfo          *gopsutil.Process
}
