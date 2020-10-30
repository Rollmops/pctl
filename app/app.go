package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Context struct {
	Config       *Config
	OutputWriter *os.File
	Processes    ProcessList
	Cache        Cache
}

var CurrentContext *Context

func init() {
	CurrentContext = &Context{}
	yamlLoader := &YamlLoader{}
	SuffixConfigLoaderMap["yaml"] = yamlLoader
	SuffixConfigLoaderMap["yml"] = yamlLoader
}

func (c *Context) GetProcessByConfig(processConfig *ProcessConfig) (*Process, error) {
	for _, p := range c.Processes {
		if p.Config.Name == processConfig.Name && p.Config.Group == processConfig.Group {
			return p, nil
		}
	}
	return nil, fmt.Errorf("unable to find process for config '%v'", processConfig)
}

func (c *Context) Initialize() error {
	logrus.Debug("Start initializing app context")
	err := c.LoadConfig()
	if err != nil {
		return err
	}
	c.RefreshProcessesFromConfig()
	err = c.Cache.Refresh()
	if err != nil {
		return err
	}
	return CurrentContext.Processes.SyncRunningInfo()

}

func (c *Context) LoadConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	logrus.Debugf("Using Config path: %s", configPath)
	configLoader := GetLoaderFromPath(configPath)
	c.Config, err = configLoader.Load(configPath)
	if err != nil {
		return err
	}
	err = c.Config.FillDependsOnInverse()
	if err != nil {
		return err
	}
	c.Config.ExpandVars()
	logrus.Debugf("Loaded %d process configuration(s)", len(c.Config.ProcessConfigs))
	return ValidateAcyclicDependencies()
}

func (c *Context) RefreshProcessesFromConfig() {
	c.Processes = make(ProcessList, 0)
	for _, processConfig := range c.Config.ProcessConfigs {
		process := &Process{Config: processConfig}
		c.Processes = append(c.Processes, process)
	}
}

func Run(args []string) error {
	pctlApp, err := CreateCliApp()
	if err != nil {
		return err
	}
	return pctlApp.Run(args)
}
