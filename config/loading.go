package config

import (
	"fmt"
	"io/ioutil"

	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type baseConfig struct {
	Include []string
}

type ProcessConfig struct {
	Name string
	Cmd  string
}

type Config struct {
	Processes []ProcessConfig
}

func LoadConfig() Config {
	log.Debug("Start loading config")
	_baseConfig := loadBaseConfig()
	log.Print(_baseConfig)

	var processConfigs []ProcessConfig

	for _, include := range _baseConfig.Include {
		log.Debug("Including config ", include)
		processConfigs = append(processConfigs, loadProcessConfigs(include)...)
	}

	config := Config{
		Processes: processConfigs,
	}

	return config
}

func loadBaseConfig() baseConfig {
	configPath := GetConfigPath()
	log.Debug("Loading basic config from ", configPath)
	data := loadFileContent(configPath)

	var _baseConfig baseConfig

	if err := yaml.Unmarshal(data, &_baseConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading YAML %s: %v\n", configPath, err)
		os.Exit(1)
	}

	return _baseConfig
}

func loadProcessConfigs(pathPattern string) []ProcessConfig {

	var processConfigs []ProcessConfig
	matches, err := filepath.Glob(ReplaceEnvVarsAndTilde(pathPattern))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading path pattern %s: %s\n", pathPattern, err)
		os.Exit(1)
	}

	for _, match := range matches {
		data := loadFileContent(match)

		var processConfig ProcessConfig
		if err := yaml.Unmarshal(data, &processConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading YAML %s: %v\n", match, err)
			os.Exit(1)
		}
		processConfigs = append(processConfigs, processConfig)

	}

	return processConfigs
}

func loadFileContent(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load file path %s\n", path)
		os.Exit(1)
	}
	return data
}
