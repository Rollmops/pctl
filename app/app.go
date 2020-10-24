package app

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Context struct {
	Config       *Config
	OutputWriter *os.File
}

var CurrentContext *Context

func init() {
	CurrentContext = &Context{}
	yamlLoader := &YamlLoader{}
	SuffixConfigLoaderMap["yaml"] = yamlLoader
	SuffixConfigLoaderMap["yml"] = yamlLoader
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
	log.Debugf("Using Config path: %s", configPath)
	configLoader := GetLoaderFromPath(configPath)
	CurrentContext.Config, err = configLoader.Load(configPath)
	if err != nil {
		return err
	}
	CurrentContext.Config.FillDependsOnInverse()
	log.Debugf("Loaded %d process configuration(s)", len(CurrentContext.Config.Processes))
	err = ValidateAcyclicDependencies()
	if err != nil {
		return err
	}

	return pctlApp.Run(args)
}
