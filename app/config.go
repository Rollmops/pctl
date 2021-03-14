package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stevenle/topsort"
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
	Command []string `yaml:"command"`
	Traits  []string `yaml:"traits"`
}

type AdditionalProcessConfig struct {
	DependsOn    []string          `yaml:"dependsOn"`
	Metadata     map[string]string `yaml:"metadata"`
	Stop         *StopConfig       `yaml:"stop"`
	StartupProbe *Probe            `yaml:"startupProbe"`
	Env          map[string]string `yaml:"env"`
	Disabled     bool              `yaml:"disabled"`
}

type ProcessConfig struct {
	CoreProcessConfig       `yaml:",inline"`
	AdditionalProcessConfig `yaml:",inline"`
	_flatPropertyMap        *map[string]string
}

type Config struct {
	PromptForStop  bool                                `yaml:"promptForStop"`
	PromptForStart bool                                `yaml:"promptForStart"`
	ProcessConfigs []*ProcessConfig                    `yaml:"processes"`
	Traits         map[string]*AdditionalProcessConfig `yaml:"traits"`
}

func (c *Config) Initialize() error {
	c.expandVars()
	err := c.mergeTraitsConfig()
	if err != nil {
		return err
	}
	return c.validateAcyclicDependencies()
}

func (c *ProcessConfig) String() string {
	if c.Group == "" {
		return c.Name
	}
	return fmt.Sprintf("%s:%s", c.Group, c.Name)
}

func (c *ProcessConfig) GetStopConfig() *StopConfig {
	if c.Stop == nil {
		return &StopConfig{}
	}
	return c.Stop
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
				p := CurrentContext.GetProcessByConfig(processConfig)
				if p == nil {
					return nil, fmt.Errorf("unable to find process config %v", processConfig)
				}
				returnProcesses = append(returnProcesses, p)
			}
		}
	}
	logrus.Tracef("Found %d process configs for name specifiers: %v", len(returnProcesses), nameSpecifiers)
	return getFilteredProcesses(returnProcesses, filters)
}

func (c *Config) GetDependencyGraph(reverse bool) (*topsort.Graph, error) {
	dependencyTree := topsort.NewGraph()
	for _, processConfig := range c.ProcessConfigs {
		dependencyTree.AddNode(processConfig.String())
		for _, dependencySpecifier := range processConfig.DependsOn {
			dependencyConfigs, err := c.FindByGroupNameSpecifier(dependencySpecifier)
			if err != nil {
				return nil, err
			}
			for _, dependencyConfig := range dependencyConfigs {
				dependencyTree.AddNode(dependencyConfig.String())
				if !reverse {
					err = dependencyTree.AddEdge(processConfig.String(), dependencyConfig.String())
				} else {
					err = dependencyTree.AddEdge(dependencyConfig.String(), processConfig.String())
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return dependencyTree, nil
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

	return "", fmt.Errorf("Unable to to find valid config Path: %v\n", possibleConfigPaths)
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

func (c *Config) mergeTraitsConfig() error {
	for _, processConfig := range c.ProcessConfigs {
		for _, trait := range processConfig.Traits {
			traitConfig := c.Traits[trait]
			if traitConfig == nil {
				return fmt.Errorf("unable to find trait '%s' for process '%s'", trait, processConfig.Name)
			}
			processConfig.Env = MergeStringMap(traitConfig.Env, processConfig.Env)
			processConfig.Metadata = MergeStringMap(traitConfig.Metadata, processConfig.Metadata)
			if processConfig.Stop == nil {
				processConfig.Stop = traitConfig.Stop
			}
			if len(processConfig.DependsOn) == 0 {
				processConfig.DependsOn = traitConfig.DependsOn
			}
			if !processConfig.Disabled {
				processConfig.Disabled = traitConfig.Disabled
			}
			if processConfig.StartupProbe == nil {
				processConfig.StartupProbe = traitConfig.StartupProbe
			}
		}
	}
	return nil
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
