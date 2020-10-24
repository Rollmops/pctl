package output

import (
	"bytes"
	"encoding/json"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
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
	Pid              int32                   `json:"pid"`
	Cwd              string                  `json:"cwd"`
	IsRunning        bool                    `json:"isRunning"`
	CPUPercent       float64                 `json:"cpuPercent"`
	Connections      []net.ConnectionStat    `json:"connections"`
	Command          []string                `json:"command"`
	MemoryInfo       *process.MemoryInfoStat `json:"memoryInfo"`
	Exe              string                  `json:"exe"`
	Username         string                  `json:"username"`
	Terminal         string                  `json:"terminal"`
	CreateTime       int64                   `json:"createTime"`
	CreateTimeString string                  `json:"createTimeString"`
}

type JsonInfoEntry struct {
	Name                 string       `json:"name"`
	ConfiguredCommand    []string     `json:"configuredCommand"`
	RunningCommand       []string     `json:"runningCommand"`
	IsRunning            bool         `json:"isRunning"`
	ConfigCommandChanged bool         `json:"configCommandChanged"`
	Info                 *ProcessInfo `json:"info"`
}

func (j *JsonConsoleOutput) SetWriter(writer *os.File) {
	j.writer = writer
}

func (j *JsonConsoleOutput) Write(infoEntries []*Info) error {

	var jsonInfoEntries []*JsonInfoEntry

	for _, infoEntry := range infoEntries {
		runningInfo, err := getRunningInfo(infoEntry)
		if err != nil {
			return err
		}
		jsonInfoEntry := &JsonInfoEntry{
			Name:                 infoEntry.Name,
			ConfiguredCommand:    infoEntry.ConfigCommand,
			RunningCommand:       infoEntry.RunningCommand,
			IsRunning:            infoEntry.IsRunning,
			Info:                 runningInfo,
			ConfigCommandChanged: infoEntry.ConfigCommandChanged,
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

func getRunningInfo(info *Info) (*ProcessInfo, error) {
	if info.RunningInfo == nil {
		return nil, nil
	}
	var runningInfo *ProcessInfo

	cwd, err := info.RunningInfo.Cwd()
	if err != nil {
		return nil, err
	}
	isRunning, err := info.RunningInfo.IsRunning()
	if err != nil {
		return nil, err
	}

	cpuPercent, err := info.RunningInfo.CPUPercent()
	if err != nil {
		return nil, err
	}

	connections, err := info.RunningInfo.Connections()
	if err != nil {
		return nil, err
	}
	command, err := info.RunningInfo.CmdlineSlice()
	if err != nil {
		return nil, err
	}
	memoryInfo, err := info.RunningInfo.MemoryInfo()
	if err != nil {
		return nil, err
	}
	exe, err := info.RunningInfo.Exe()
	if err != nil {
		return nil, err
	}
	createTime, err := info.RunningInfo.CreateTime()
	if err != nil {
		return nil, err
	}
	username, err := info.RunningInfo.Username()
	if err != nil {
		return nil, err
	}
	terminal, err := info.RunningInfo.Terminal()
	if err != nil {
		return nil, err
	}
	createTimeString := time.Unix(createTime/1000, 0)
	if info.IsRunning {
		runningInfo = &ProcessInfo{
			Pid:              info.RunningInfo.Pid,
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
	}
	return runningInfo, nil
}
