package app

import (
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/output"
	"github.com/Rollmops/pctl/persistence"
	log "github.com/sirupsen/logrus"
	"os"
)

type Context struct {
	config            *config.Config
	persistenceWriter persistence.Writer
	persistenceReader persistence.Reader
	output            output.Output
}

var CurrentContext *Context

func NewContext() (*Context, error) {
	persistenceWriter, err := persistence.NewCsvWriter()
	if err != nil {
		return nil, err
	}
	persistenceReader, err := persistence.NewCsvReader()
	if err != nil {
		return nil, err
	}
	return &Context{
		persistenceWriter: persistenceWriter,
		persistenceReader: persistenceReader,
		output:            output.NewSimpleConsoleOutput(os.Stdout),
	}, nil
}

func (c *Context) Initialize() error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}
	log.Debugf("Using config path: %s", configPath)
	configLoader := config.GetLoaderFromPath(configPath)
	c.config, err = configLoader.Load(configPath)
	log.Debugf("Loaded %d process configuration(s)", len(c.config.Processes))
	if err != nil {
		return err
	}
	return nil
}
