package app

import (
	"bytes"
	"encoding/json"
	"github.com/shirou/gopsutil/net"
	gopsutil "github.com/shirou/gopsutil/process"
	"io"
	"os"
	"time"
)

func init() {
	FormatMap["json"] = &JsonConsoleOutput{
		indent: "  ",
		flat:   false,
	}
	FormatMap["json-flat"] = &JsonConsoleOutput{
		indent: "",
		flat:   true,
	}
}

type JsonConsoleOutput struct {
	writer io.Writer
	indent string
	flat   bool
}

type ProcessInfo struct {
	Pid              int32                    `json:"pid"`
	Nice             int32                    `json:"nice"`
	Cwd              string                   `json:"cwd"`
	IsRunning        bool                     `json:"isRunning"`
	CPUPercent       float64                  `json:"cpuPercent"`
	Connections      []net.ConnectionStat     `json:"connections"`
	Command          []string                 `json:"command"`
	MemoryInfo       *gopsutil.MemoryInfoStat `json:"memoryInfo"`
	Exe              string                   `json:"exe"`
	Username         string                   `json:"username"`
	Terminal         string                   `json:"terminal"`
	CreateTime       int64                    `json:"createTime"`
	CreateTimeString string                   `json:"createTimeString"`
}

type JsonInfoEntry struct {
	Name              string       `json:"name"`
	Group             string       `json:"group"`
	ConfiguredCommand []string     `json:"configuredCommand"`
	RunningCommand    []string     `json:"runningCommand"`
	IsRunning         bool         `json:"isRunning"`
	DirtyCommand      bool         `json:"dirtyCommand"`
	DirtyCommandArgs  []string     `json:"dirtyCommandArgs"`
	Dirty             bool         `json:"dirty"`
	Info              *ProcessInfo `json:"info"`
}

func (j *JsonConsoleOutput) SetWriter(writer *os.File) {
	j.writer = writer
}

func (j *JsonConsoleOutput) Write(processes ProcessList, _ []string) error {
	var jsonInfoEntries []*JsonInfoEntry

	for _, p := range processes {
		jsonInfoEntry := &JsonInfoEntry{
			Name:              p.Config.Name,
			ConfiguredCommand: p.Config.Command,
		}
		var jsonRunningInfo *ProcessInfo
		if p.IsRunning() {
			var err error
			jsonRunningInfo, err = getJsonProcessInfo(p.Info.GoPsutilProcess)
			if err != nil {
				return err
			}
			jsonInfoEntry.IsRunning = true
			jsonInfoEntry.RunningCommand = p.Info.RunningCommand
			jsonInfoEntry.Info = jsonRunningInfo
			jsonInfoEntry.Dirty = p.Info.Dirty
			jsonInfoEntry.DirtyCommand = p.Info.DirtyCommand
			jsonInfoEntry.DirtyCommandArgs = p.Info.DirtyMd5Hashes

		}
		jsonInfoEntries = append(jsonInfoEntries, jsonInfoEntry)
	}

	b, err := json.Marshal(jsonInfoEntries)
	if err != nil {
		return err
	}
	b = append(b, []byte("\n")...)

	var outBuffer bytes.Buffer

	if !j.flat {
		err = json.Indent(&outBuffer, b, "", j.indent)
		if err != nil {
			return err
		}
		_, err = j.writer.Write(outBuffer.Bytes())

		return err
	} else {
		_, err = j.writer.Write(b)
		return err
	}
}

func getJsonProcessInfo(gopsutilProcess *gopsutil.Process) (*ProcessInfo, error) {
	var runningInfo *ProcessInfo

	cwd, err := gopsutilProcess.Cwd()
	if err != nil {
		return nil, err
	}
	isRunning, err := gopsutilProcess.IsRunning()
	if err != nil {
		return nil, err
	}

	cpuPercent, err := gopsutilProcess.CPUPercent()
	if err != nil {
		return nil, err
	}

	connections, err := gopsutilProcess.Connections()
	if err != nil {
		return nil, err
	}
	command, err := gopsutilProcess.CmdlineSlice()
	if err != nil {
		return nil, err
	}
	memoryInfo, err := gopsutilProcess.MemoryInfo()
	if err != nil {
		return nil, err
	}
	exe, err := gopsutilProcess.Exe()
	if err != nil {
		return nil, err
	}
	createTime, err := gopsutilProcess.CreateTime()
	if err != nil {
		return nil, err
	}
	username, err := gopsutilProcess.Username()
	if err != nil {
		return nil, err
	}
	terminal, err := gopsutilProcess.Terminal()
	if err != nil {
		return nil, err
	}
	createTimeString := time.Unix(createTime/1000, 0)
	nice, err := gopsutilProcess.Nice()
	if err != nil {
		return nil, err
	}
	runningInfo = &ProcessInfo{
		Nice:             nice,
		Pid:              gopsutilProcess.Pid,
		Cwd:              cwd,
		IsRunning:        isRunning,
		CPUPercent:       cpuPercent,
		Connections:      connections,
		Command:          command,
		MemoryInfo:       memoryInfo,
		Exe:              exe,
		Username:         username,
		Terminal:         terminal,
		CreateTime:       createTime,
		CreateTimeString: createTimeString.String(),
	}
	return runningInfo, nil
}
