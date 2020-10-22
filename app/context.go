package app

import (
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/config/yaml"
	log "github.com/sirupsen/logrus"
	"os"
)

type Context struct {
	Config       *config.Config
	OutputWriter *os.File
}

var CurrentContext *Context

func init() {
	CurrentContext = &Context{}
	loader := &yaml.Loader{}
	config.SuffixConfigLoaderMap["yaml"] = loader
	config.SuffixConfigLoaderMap["yml"] = loader
}

func (c *Context) Initialize() error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}
	log.Debugf("Using Config path: %s", configPath)
	configLoader := config.GetLoaderFromPath(configPath)
	c.Config, err = configLoader.Load(configPath)
	c.Config.FillDependsOnInverse()
	if err != nil {
		return err
	}
	log.Debugf("Loaded %d process configuration(s)", len(c.Config.Processes))
	return nil
}
