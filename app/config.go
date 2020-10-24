package app

import (
	"fmt"
	"github.com/minio/minio/pkg/wildcard"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"syscall"
)

const _configFileName string = "pctl.yml"

type Loader interface {
	Load(path string) (*Config, error)
}

type StopStrategyConfig struct {
	Script       *ScriptStopStrategyConfig
	Signal       *SignalStopStrategyConfig
	MaxWaitTime  string
	IntervalTime string
}

type SignalStopStrategyConfig struct {
	Signal       syscall.Signal
	SignalString string
}

type ScriptStopStrategyConfig struct {
	Path          string
	Args          []string
	ForwardStdout bool `yaml:"forwardStdout"`
	ForwardStderr bool `yaml:"forwardStderr"`
}

type ProcessConfig struct {
	Name                    string
	Command                 []string            `yaml:"cmd"`
	PidRetrieveStrategyName string              `yaml:"pidStrategy"`
	StopStrategy            *StopStrategyConfig `yaml:"stop"`
	DependsOn               []string            `yaml:"dependsOn"`
	DependsOnInverse        []string
	Metadata                map[string]string `yaml:"metadata"`
}

type Config struct {
	Processes []*ProcessConfig
}

var SuffixConfigLoaderMap = map[string]Loader{}

func (c *Config) FindByName(name string) *ProcessConfig {
	log.Tracef("Getting process config for name '%s'", name)
	for _, p := range c.Processes {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (c *Config) CollectSyncedProcessesByNameSpecifiers(nameSpecifiers []string, filters []string, allIfNoSpecifiers bool) ([]*Process, error) {
	log.Tracef("Collecting processes for name specifiers: %v", nameSpecifiers)
	var returnProcesses []*Process
	if len(nameSpecifiers) == 0 && allIfNoSpecifiers {
		for _, pConfig := range c.Processes {
			p, err := NewFromConfigAndSynced(pConfig)
			if err != nil {
				return nil, err
			}
			returnProcesses = append(returnProcesses, p)
		}
	}
	for _, nameSpecifier := range nameSpecifiers {
		for _, processConfig := range c.Processes {
			if wildcard.Match(nameSpecifier, processConfig.Name) && !_isInProcessList(processConfig.Name, returnProcesses) {
				p, err := NewFromConfigAndSynced(processConfig)
				if err != nil {
					return nil, err
				}
				returnProcesses = append(returnProcesses, p)
			}
		}
	}
	log.Tracef("Found %d process configs for name specifiers: %v", len(returnProcesses), nameSpecifiers)
	return getFilteredProcesses(returnProcesses, filters)
}

func getFilteredProcesses(processes []*Process, filterPatterns []string) ([]*Process, error) {
	if len(filterPatterns) > 0 {
		var filteredProcesses []*Process
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
	for _, pConfig := range c.Processes {
		for _, dependsOn := range pConfig.DependsOn {
			dConfig := c.FindByName(dependsOn)
			if !_isInList(dConfig.DependsOnInverse, pConfig.Name) {
				dConfig.DependsOnInverse = append(dConfig.DependsOnInverse, pConfig.Name)
			}
		}
	}
}

func _isInProcessList(name string, processes []*Process) bool {
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
