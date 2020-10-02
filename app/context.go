package app

import (
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/config/yaml"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/persistence/csv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	loader := &yaml.Loader{}
	config.SuffixConfigLoaderMap["yaml"] = loader
	config.SuffixConfigLoaderMap["yml"] = loader

	log.Trace("Creating context")
	persistenceWriter, err := csv.NewCsvWriter()
	if err != nil {
		log.Fatal(err.Error())
	}
	persistenceReader, err := csv.NewCsvReader()
	if err != nil {
		log.Fatal(err.Error())
	}
	CurrentContext = &Context{
		PersistenceWriter: persistenceWriter,
		PersistenceReader: persistenceReader,
	}
}

type Context struct {
	Config            *config.Config
	PersistenceWriter persistence.Writer
	PersistenceReader persistence.Reader
	OutputWriter      *os.File
}

var CurrentContext *Context

func (c *Context) Initialize() error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}
	log.Debugf("Using Config path: %s", configPath)
	configLoader := config.GetLoaderFromPath(configPath)
	c.Config, err = configLoader.Load(configPath)
	if err != nil {
		return err
	}
	log.Debugf("Loaded %d process configuration(s)", len(c.Config.Processes))
	return nil
}
