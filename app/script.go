package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

type Script struct {
	Path          string
	Args          []string
	ForwardStdout bool `yaml:"forwardStdout"`
	ForwardStderr bool `yaml:"forwardStderr"`
}

func (s *Script) CreateCommand(process *Process) (*exec.Cmd, error) {
	mapping := process.Config.GetFlatPropertyMap()
	if process.RunningInfo != nil {
		mapping["pid"] = strconv.Itoa(int(process.RunningInfo.Pid))
	}
	scriptPath, err := filepath.Abs(ExpandPath(s.Path))
	if err != nil {
		return nil, err
	}

	var substArgs []string
	for _, arg := range s.Args {
		substArg := replaceVariables(arg, mapping)
		substArgs = append(substArgs, substArg)
	}
	cmd := exec.Command(scriptPath, substArgs...)
	if s.ForwardStdout {
		cmd.Stdout = os.Stdout
	}
	if s.ForwardStderr {
		cmd.Stderr = os.Stderr
	}
	return cmd, nil
}

func replaceVariables(text string, replaceMap map[string]string) string {
	re := regexp.MustCompile(`\$\{(\S+)\}`)
	s := re.ReplaceAllStringFunc(text, func(s string) string {
		matches := re.FindAllStringSubmatch(s, 1)
		if len(matches) == 1 {
			return replaceMap[matches[0][1]]
		}
		return text
	})
	return s
}
