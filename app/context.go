package app

import (
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/persistence"
)

type Context struct {
	config            *config.Config
	persistenceWriter persistence.Writer
	persistenceReader persistence.Reader
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
	}, nil
}

func (c *Context) Initialize() error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}
	configLoader := config.GetLoaderFromPath(configPath)
	c.config, err = configLoader.Load(configPath)
	if err != nil {
		return err
	}
	return nil
}
