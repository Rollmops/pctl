package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Context struct {
	Config           *Config
	OutputWriter     *os.File
	RunningProcesses []*Process
}

var CurrentContext *Context

func init() {
	CurrentContext = &Context{}
	yamlLoader := &YamlLoader{}
	SuffixConfigLoaderMap["yaml"] = yamlLoader
	SuffixConfigLoaderMap["yml"] = yamlLoader
}

func (c *Context) GetProcessByName(name string) (*Process, error) {
	for _, p := range c.RunningProcesses {
		if p.Config.Name == name {
			return p, nil
		}
	}
	config := c.Config.FindByName(name)
	if config == nil {
		return nil, fmt.Errorf("unable to find config '%s'", name)
	}
	p := &Process{Config: config}
	return p, nil
}

func Run(args []string) error {
	pctlApp, err := CreateCliApp()
	if err != nil {
		return err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	logrus.Debugf("Using Config path: %s", configPath)
	configLoader := GetLoaderFromPath(configPath)
	CurrentContext.Config, err = configLoader.Load(configPath)
	if err != nil {
		return err
	}
	CurrentContext.Config.FillDependsOnInverse()
	logrus.Debugf("Loaded %d process configuration(s)", len(CurrentContext.Config.ProcessConfigs))
	err = ValidateAcyclicDependencies()
	if err != nil {
		return err
	}

	err = CurrentContext.SyncRunningProcesses()
	if err != nil {
		return err
	}
	return pctlApp.Run(args)
}
