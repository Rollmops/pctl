package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

func init() {
	FormatMap["json"] = &JsonConsoleOutput{
		writer: os.Stdout,
		indent: "  ",
		flat:   false,
	}
	FormatMap["json_flat"] = &JsonConsoleOutput{
		writer: os.Stdout,
		indent: "",
		flat:   true,
	}
}

type JsonConsoleOutput struct {
	writer io.Writer
	indent string
	flat   bool
}

type RunningInfo struct {
	Pid        int32   `json:"pid"`
	Cwd        string  `json:"cwd"`
	IsRunning  bool    `json:"isRunning"`
	CPUPercent float64 `json:"cpuPercent"`
}

type JsonInfoEntry struct {
	Name                 string       `json:"name"`
	ConfiguredCommand    []string     `json:"configuredCommand"`
	RunningCommand       []string     `json:"runningCommand"`
	IsRunning            bool         `json:"isRunning"`
	StoppedUnexpectedly  bool         `json:"stoppedUnexpectedly"`
	ConfigCommandChanged bool         `json:"configCommandChanged"`
	Info                 *RunningInfo `json:"info"`
}

func (j *JsonConsoleOutput) Write(infoEntries []*InfoEntry) error {

	var jsonInfoEntries []*JsonInfoEntry

	for _, infoEntry := range infoEntries {
		runningInfo, err := getRunningInfo(infoEntry)
		if err != nil {
			return err
		}
		jsonInfoEntry := &JsonInfoEntry{
			Name:              infoEntry.Name,
			ConfiguredCommand: infoEntry.ConfigCommand,
			RunningCommand:    infoEntry.RunningCommand,
			Info:              runningInfo,
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

func getRunningInfo(infoEntry *InfoEntry) (*RunningInfo, error) {
	if infoEntry.RunningInfo == nil {
		return nil, nil
	}
	var runningInfo *RunningInfo

	cwd, err := infoEntry.RunningInfo.Cwd()
	if err != nil {
		return nil, err
	}
	isRunning, err := infoEntry.RunningInfo.IsRunning()
	if err != nil {
		return nil, err
	}

	cpuPercent, err := infoEntry.RunningInfo.CPUPercent()
	if err != nil {
		return nil, err
	}
	if infoEntry.IsRunning {
		runningInfo = &RunningInfo{
			Pid:        infoEntry.RunningInfo.Pid,
			Cwd:        cwd,
			IsRunning:  isRunning,
			CPUPercent: cpuPercent,
		}
	}
	return runningInfo, nil
}
