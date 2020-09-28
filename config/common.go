package config

import (
	"fmt"
	"github.com/minio/minio/pkg/wildcard"
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

func (c *Config) CollectProcessConfigsByNameSpecifiers(nameSpecifiers []string, allIfNoSpecifiers bool) []*ProcessConfig {
	if len(nameSpecifiers) == 0 && allIfNoSpecifiers {
		return c.Processes
	}
	log.Tracef("Collecting process configs for name specifiers: %v", nameSpecifiers)
	var returnConfigs []*ProcessConfig
	for _, nameSpecifier := range nameSpecifiers {
		for _, p := range c.Processes {
			if wildcard.Match(nameSpecifier, p.Name) && !_isInProcessConfigList(p.Name, returnConfigs) {
				returnConfigs = append(returnConfigs, p)
			}
		}
	}
	log.Tracef("Found %d process configs for name specifiers: %v", len(returnConfigs), nameSpecifiers)
	return returnConfigs
}

func (c *Config) GetAllDependentOf(name string) []*ProcessConfig {
	var dependentReturns []*ProcessConfig
	for _, p := range c.Processes {
		for _, d := range p.DependsOn {
			if d == name && !_isInProcessConfigList(name, dependentReturns) {
				dependentReturns = append(dependentReturns, p)
			}
		}
	}
	return dependentReturns
}

func _isInProcessConfigList(name string, processConfigs []*ProcessConfig) bool {
	for _, p := range processConfigs {
		if name == p.Name {
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

	return "", fmt.Errorf("Unable to to find valid config path: %v\n", possibleConfigPaths)
}

func GetLoaderFromPath(path string) Loader {
	fractions := strings.Split(path, ".")
	suffix := fractions[len(fractions)-1]
	return SuffixConfigLoaderMap[suffix]
}
