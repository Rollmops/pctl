package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

const _configFileName string = "pctl.yml"

type Loader interface {
	Load(path string) (*Config, error)
}

type ProcessConfig struct {
	Name             string
	Group            string
	Command          []string            `yaml:"cmd"`
	StopStrategy     *StopStrategyConfig `yaml:"stop"`
	DependsOn        []string            `yaml:"dependsOn"`
	DependsOnInverse []string
	Metadata         map[string]string `yaml:"metadata"`
	ReadinessProbe   string            `yaml:"readinessProbe"`
	StartupProbe     string            `yaml:"startupProbe"`
	Env              map[string]string `yaml:"env"`
}

type Config struct {
	ProcessConfigs []*ProcessConfig
}

func (c *Config) ExpandVars() {
	for _, pConfig := range c.ProcessConfigs {
		var command []string
		for _, arg := range pConfig.Command {
			command = append(command, ExpandPath(arg))
		}
		pConfig.Command = command
	}
}

var SuffixConfigLoaderMap = make(map[string]Loader)

func (c *ProcessConfig) String() string {
	return fmt.Sprintf("ProcessConfig(group: %s, name: %s)", c.Group, c.Name)
}

func (c *Config) FindByName(name string) *ProcessConfig {
	logrus.Tracef("Getting process config for name '%s'", name)
	for _, p := range c.ProcessConfigs {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (c *Config) CollectProcessesByNameSpecifiers(nameSpecifiers []string, filters []string, allIfNoSpecifiers bool) (ProcessList, error) {
	logrus.Tracef("Collecting processes for name specifiers: %v", nameSpecifiers)
	var returnProcesses ProcessList
	if len(nameSpecifiers) == 0 && allIfNoSpecifiers {
		returnProcesses = CurrentContext.Processes
	}
	for _, specifier := range nameSpecifiers {
		groupNameSpecifier, err := NewGroupNameSpecifier(specifier)
		if err != nil {
			return nil, err
		}
		for _, processConfig := range c.ProcessConfigs {
			if groupNameSpecifier.IsMatchingGroupAndName(processConfig.Group, processConfig.Name) && !_isInProcessList(processConfig.Name, returnProcesses) {
				p, err := CurrentContext.GetProcessByConfig(processConfig)
				if err != nil {
					return nil, err
				}
				returnProcesses = append(returnProcesses, p)
			}
		}
	}
	logrus.Tracef("Found %d process configs for name specifiers: %v", len(returnProcesses), nameSpecifiers)
	return getFilteredProcesses(returnProcesses, filters)
}

func getFilteredProcesses(processes ProcessList, filterPatterns []string) ([]*Process, error) {
	err := processes.SyncRunningInfo()
	if err != nil {
		return nil, err
	}
	if len(filterPatterns) > 0 {
		var filteredProcesses ProcessList
		for _, filterPattern := range filterPatterns {
			filter, err := NewFilter(filterPattern)
			if err != nil {
				return nil, err
			}
			for _, p := range processes {
				isRelevant, err := filter.IsMatchingProcess(p)
				if err != nil {
					return nil, err
				}
				if isRelevant {
					filteredProcesses = append(filteredProcesses, p)
				}
			}
		}
		return filteredProcesses, nil
	}
	return processes, nil
}

func (c *Config) FillDependsOnInverse() {
	for _, pConfig := range c.ProcessConfigs {
		for _, dependsOn := range pConfig.DependsOn {
			dConfig := c.FindByName(dependsOn)
			if !_isInList(dConfig.DependsOnInverse, pConfig.Name) {
				dConfig.DependsOnInverse = append(dConfig.DependsOnInverse, pConfig.Name)
			}
		}
	}
}

func _isInProcessList(name string, processes ProcessList) bool {
	for _, p := range processes {
		if name == p.Config.Name {
			return true
		}
	}
	return false
}

func GetConfigPath() (string, error) {
	cwd, _ := os.Getwd()
	possibleConfigPaths := []string{
		os.Getenv("PCTL_CONFIG_PATH"),
		path.Join(cwd, _configFileName),
		path.Join(os.Getenv("HOME"), ".config", _configFileName),
		path.Join("/", "etc", "pctl", _configFileName),
	}

	for _, configPath := range possibleConfigPaths {
		_, err := os.Stat(configPath)
		if err == nil {
			return configPath, nil
		}
	}

	return "", fmt.Errorf("Unable to to find valid config path: %v\n", possibleConfigPaths)
}

func GetLoaderFromPath(path string) Loader {
	fractions := strings.Split(path, ".")
	suffix := fractions[len(fractions)-1]
	return SuffixConfigLoaderMap[suffix]
}
