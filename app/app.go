package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Context struct {
	Config       *Config
	OutputWriter *os.File
	Processes    ProcessList
	Cache        Cache
	Context      string
}

var CurrentContext *Context

func init() {
	CurrentContext = &Context{}
	yamlLoader := &YamlLoader{}
	SuffixConfigLoaderMap["yaml"] = yamlLoader
	SuffixConfigLoaderMap["yml"] = yamlLoader
}

func (c *Context) GetProcessByConfig(processConfig *ProcessConfig) *Process {
	for _, p := range c.Processes {
		if p.Config.Name == processConfig.Name && p.Config.Group == processConfig.Group {
			return p
		}
	}
	return nil
}

func (c *Context) GetProcessEnvironmentMarker() string {
	return fmt.Sprintf("__PCTL_INFO_%s__", strings.ToUpper(c.Context))
}

func (c *Context) InitializeRunningProcessInfo() error {
	logrus.Debug("Start initializing app context")
	c.RefreshProcessesFromConfig()
	err := c.Cache.Refresh()
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
	logrus.Debugf("Using Config Path: %s", configPath)
	configLoader := GetLoaderFromPath(configPath)
	c.Config, err = configLoader.Load(configPath)
	if err != nil {
		return err
	}
	return c.Config.Initialize()
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
