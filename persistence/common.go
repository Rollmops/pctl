package persistence

import (
	"fmt"
	"github.com/Rollmops/pctl/process"
	"os"
	"path/filepath"
)

func GetStateFilePath() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("unable to retrieve HOME env var")
	}
	stateFilePath := filepath.Join(home, "var", "pctl")

	if _, err := os.Stat(stateFilePath); err != nil {
		err = os.MkdirAll(stateFilePath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return stateFilePath, nil
}

func NewDataEntryFromProcess(p *process.Process) (*DataEntry, error) {
	pid, err := p.Pid()
	if err != nil {
		return nil, err
	}
	info, err := p.Info()
	if err != nil {
		return nil, err
	}
	cmdline, err := info.Cmdline()
	if err != nil {
		return nil, err
	}
	return &DataEntry{
		Pid:  pid,
		Name: p.Config.Name,
		Cmd:  cmdline,
	}, nil
}
