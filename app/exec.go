package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

type Exec struct {
	Command       []string `yaml:"command"`
	ForwardStdout bool     `yaml:"forwardStdout"`
	ForwardStderr bool     `yaml:"forwardStderr"`
}

func (s *Exec) CreateCommand(process *Process) (*exec.Cmd, error) {
	if len(s.Command) == 0 {
		return nil, fmt.Errorf("command length of exec probe is 0")
	}
	mapping := process.Config.GetFlatPropertyMap()
	if process.RunningInfo != nil {
		mapping["pid"] = strconv.Itoa(int(process.RunningInfo.Pid))
	}

	name, err := filepath.Abs(ExpandPath(s.Command[0]))
	if err != nil {
		return nil, err
	}
	var args []string
	if len(s.Command) > 1 {
		for _, arg := range s.Command[1:] {
			substArg := replaceVariables(arg, mapping)
			args = append(args, substArg)
		}
	}
	cmd := exec.Command(name, args...)
	if s.ForwardStdout {
		cmd.Stdout = os.Stdout
	}
	if s.ForwardStderr {
		cmd.Stderr = os.Stderr
	}
	return cmd, nil
}

func replaceVariables(text string, replaceMap map[string]string) string {
	re := regexp.MustCompile(`\${(\S+)}`)
	s := re.ReplaceAllStringFunc(text, func(s string) string {
		matches := re.FindAllStringSubmatch(s, 1)
		if len(matches) == 1 {
			return replaceMap[matches[0][1]]
		}
		return text
	})
	return s
}
