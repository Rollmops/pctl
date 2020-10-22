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

func (c *Config) CollectProcessConfigsByNameSpecifiers(nameSpecifiers []string, filters []string, allIfNoSpecifiers bool) ([]*ProcessConfig, error) {
	log.Tracef("Collecting process configs for name specifiers: %v", nameSpecifiers)
	var returnConfigs []*ProcessConfig
	if len(nameSpecifiers) == 0 && allIfNoSpecifiers {
		returnConfigs = c.Processes
	}
	for _, nameSpecifier := range nameSpecifiers {
		for _, p := range c.Processes {
			if wildcard.Match(nameSpecifier, p.Name) && !_isInProcessConfigList(p.Name, returnConfigs) {
				returnConfigs = append(returnConfigs, p)
			}
		}
	}
	log.Tracef("Found %d process configs for name specifiers: %v", len(returnConfigs), nameSpecifiers)
	return getFilteredProcessConfigs(returnConfigs, filters)
}

func getFilteredProcessConfigs(processConfigs []*ProcessConfig, filters []string) ([]*ProcessConfig, error) {
	if len(filters) > 0 {
		var filteredProcessConfigs []*ProcessConfig
		for _, filter := range filters {
			fractions := strings.Split(filter, "=")
			if len(fractions) != 2 {
				return nil, fmt.Errorf("invalid filter format: '%s'", filter)
			}
			for _, pConfig := range processConfigs {
				if fractions[0] == "label" {
					if _isInList(pConfig.Labels, fractions[1]) {
						filteredProcessConfigs = append(filteredProcessConfigs, pConfig)
					}
				} else {
					value := pConfig.Metadata[fractions[0]]
					if value == fractions[1] {
						filteredProcessConfigs = append(filteredProcessConfigs, pConfig)
					}
				}
			}
		}
		return filteredProcessConfigs, nil
	}
	return processConfigs, nil
}

func (c *Config) FillDependsOnInverse() {
	for _, pConfig := range c.Processes {
		for _, dependsOn := range pConfig.DependsOn {
			dConfig := c.FindByName(dependsOn)
			if !_isInList(dConfig.DependsOnInverse, pConfig.Name) {
				dConfig.DependsOnInverse = append(dConfig.DependsOnInverse, pConfig.Name)
			}
		}
	}
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
