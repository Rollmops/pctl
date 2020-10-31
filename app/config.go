package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yourbasic/graph"
	"os"
	"path"
	"strings"
)

const _configFileName string = "pctl.yml"

var SuffixConfigLoaderMap = make(map[string]Loader)

type Loader interface {
	Load(path string) (*Config, error)
}

type CoreProcessConfig struct {
	Name    string   `yaml:"name"`
	Group   string   `yaml:"group"`
	Command []string `yaml:"cmd"`
}

type AdditionalProcessConfig struct {
	WaitAfterStart   string              `yaml:"waitAfterStart"`
	StopStrategy     *StopStrategyConfig `yaml:"stop"`
	DependsOn        []string            `yaml:"dependsOn"`
	DependsOnInverse []string
	Metadata         map[string]string `yaml:"metadata"`
	ReadinessProbe   *ProbeConfig      `yaml:"readinessProbe"`
	LivenessProbe    *ProbeConfig      `yaml:"livenessProbe"`
	Env              map[string]string `yaml:"env"`
	Disabled         bool              `yaml:"disabled"`
}

type ProcessConfig struct {
	CoreProcessConfig       `yaml:",inline"`
	AdditionalProcessConfig `yaml:",inline"`
	_flatPropertyMap        *map[string]string
}

type Config struct {
	ProcessConfigs []*ProcessConfig                    `yaml:"processes"`
	Groups         map[string]*AdditionalProcessConfig `yaml:"groups"`
}

func (c *Config) Initialize() error {
	c.expandVars()
	c.mergeGroupConfig()

	err := c.validateAcyclicDependencies()
	if err != nil {
		return err
	}
	return c.fillDependsOnInverse()
}

func (c *ProcessConfig) String() string {
	if c.Group == "" {
		return c.Name
	}
	return fmt.Sprintf("%s:%s", c.Group, c.Name)
}

func (c *Config) FindByGroupAndName(group string, name string) *ProcessConfig {
	logrus.Tracef("Getting process config for name %s", name)
	for _, p := range c.ProcessConfigs {
		if p.Group == group && p.Name == name {
			return p
		}
	}
	return nil
}

func (c *Config) FindByGroupNameSpecifier(groupNameSpecifier string) ([]*ProcessConfig, error) {
	specifier, err := NewGroupNameSpecifier(groupNameSpecifier)
	if err != nil {
		return nil, err
	}
	var processConfigs []*ProcessConfig
	for _, p := range c.ProcessConfigs {
		if specifier.IsMatchingGroupAndName(p.Group, p.Name) {
			processConfigs = append(processConfigs, p)
		}
	}
	return processConfigs, nil
}

func (c *Config) CollectProcessesByNameSpecifiers(nameSpecifiers []string, filters Filters, allIfNoSpecifiers bool) (ProcessList, error) {
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

func getFilteredProcesses(processes ProcessList, filters Filters) ([]*Process, error) {
	err := processes.SyncRunningInfo()
	if err != nil {
		return nil, err
	}
	if len(filters) > 0 {
		var filteredProcesses ProcessList
		for _, filter := range filters {
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

func (c *Config) fillDependsOnInverse() error {
	for _, pConfig := range c.ProcessConfigs {
		for _, dependsOn := range pConfig.DependsOn {
			dependencyConfigs, err := c.FindByGroupNameSpecifier(dependsOn)
			if err != nil {
				return err
			}
			if dependencyConfigs == nil {
				return fmt.Errorf("unable to find process config matching %s", dependsOn)
			}
			for _, dependencyConfig := range dependencyConfigs {
				if !IsInStringList(dependencyConfig.DependsOnInverse, pConfig.Name) {
					dependencyConfig.DependsOnInverse = append(dependencyConfig.DependsOnInverse, pConfig.Name)
				}
			}
		}
	}
	return nil
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

func (c *Config) validateAcyclicDependencies() error {
	mapping := make(map[string]int)
	for i, _config := range c.ProcessConfigs {
		mapping[_config.Name] = i
	}

	gm := graph.New(len(c.ProcessConfigs))
	for _, _config := range c.ProcessConfigs {
		for _, n := range _config.DependsOn {
			gm.Add(mapping[_config.Name], mapping[n])
		}
	}
	if !graph.Acyclic(gm) {
		return fmt.Errorf("process dependency configuration is not acyclic")
	}
	return nil
}

func (c *Config) mergeGroupConfig() {
	for _, processConfig := range c.ProcessConfigs {
		groupConfig := c.Groups[processConfig.Group]
		if groupConfig != nil {
			processConfig.Env = mergeStringMap(processConfig.Env, groupConfig.Env)
			processConfig.Metadata = mergeStringMap(processConfig.Metadata, groupConfig.Metadata)
			if processConfig.WaitAfterStart == "" {
				processConfig.WaitAfterStart = groupConfig.WaitAfterStart
			}
			if processConfig.StopStrategy == nil {
				processConfig.StopStrategy = groupConfig.StopStrategy
			}
			if len(processConfig.DependsOn) == 0 {
				processConfig.DependsOn = groupConfig.DependsOn
			}
			if !processConfig.Disabled {
				processConfig.Disabled = groupConfig.Disabled
			}
			if processConfig.ReadinessProbe == nil {
				processConfig.ReadinessProbe = groupConfig.ReadinessProbe
			}
			if processConfig.LivenessProbe == nil {
				processConfig.LivenessProbe = groupConfig.LivenessProbe
			}
		}
	}
}

func mergeStringMap(processMap map[string]string, groupMap map[string]string) map[string]string {
	returnMap := make(map[string]string)
	for k, v := range groupMap {
		returnMap[k] = v
	}
	for k, v := range processMap {
		returnMap[k] = v
	}
	return returnMap
}

func (c *Config) expandVars() {
	for _, pConfig := range c.ProcessConfigs {
		var command []string
		for _, arg := range pConfig.Command {
			command = append(command, ExpandPath(arg))
		}
		pConfig.Command = command
	}
}

func (c *ProcessConfig) GetFlatPropertyMap() map[string]string {
	if c._flatPropertyMap == nil {
		flatPropertyMap := make(map[string]string)
		for k, v := range c.Metadata {
			flatPropertyMap["metadata."+k] = v
		}
		for k, v := range c.Env {
			flatPropertyMap["env."+k] = v
		}
		flatPropertyMap["name"] = c.Name
		flatPropertyMap["group"] = c.Group
		c._flatPropertyMap = &flatPropertyMap
	}
	return *c._flatPropertyMap
}
