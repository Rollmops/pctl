package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"os"

	"gopkg.in/yaml.v2"
)

type ProcessConfig struct {
	Name string
	Cmd  string
}

type Config struct {
	Processes []ProcessConfig
}

type rawConfig struct {
	Includes  []string
	Processes []ProcessConfig
}

func LoadConfig(path string) (*Config, error) {
	path, _ = filepath.Abs(os.ExpandEnv(path))
	rawConfig, err := loadYamlFromPath(path)
	if err != nil {
		return nil, err
	}

	config := Config{
		Processes: rawConfig.Processes,
	}

	err = loadIncludes(path, rawConfig.Includes, &config)
	if err != nil {
		return nil, err
	}
	err = config.validate()
	return &config, err
}

func loadIncludes(baseConfigPath string, includes []string, config *Config) error {
	for _, include := range includes {
		include, _ = filepath.Abs(os.ExpandEnv(include))

		includeMatches, err := filepath.Glob(include)
		if err != nil {
			return err
		}
		for _, include := range includeMatches {
			if baseConfigPath == include {
				return fmt.Errorf("Config file %s is trying to include itself", baseConfigPath)
			}
			rawConfig, err := loadYamlFromPath(include)
			if err != nil {
				return err
			}
			config.Processes = append(config.Processes, rawConfig.Processes...)
		}
	}
	return nil
}

func loadYamlFromPath(path string) (*rawConfig, error) {
	content, err := loadFileContent(path)
	if err != nil {
		return nil, err
	}
	var rawConfig rawConfig
	if err := yaml.Unmarshal(content, &rawConfig); err != nil {
		return nil, fmt.Errorf("Error reading YAML %s: %v", path, err)
	}

	return &rawConfig, nil
}

func loadFileContent(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to load file path %s", path)
	}
	return data, nil
}
