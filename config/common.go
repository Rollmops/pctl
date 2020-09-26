package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

const _configFileName string = "pctl.yml"

var SuffixConfigLoaderMap = map[string]Loader{}

func (c *Config) FindByName(name string) *ProcessConfig {
	log.Tracef("Getting process config for name '%s'", name)
	for _, p := range c.Processes {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (c *Config) GetAllProcessNames() []string {
	var names []string
	for _, c := range c.Processes {
		names = append(names, c.Name)
	}
	return names
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
