package config

import (
	"fmt"
	"io/ioutil"

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
	Include   []string
	Processes []ProcessConfig
}

func LoadConfig(path string) (Config, error) {
	rawConfig := loadYamlFromPath(path)

	config := Config{
		Processes: rawConfig.Processes,
	}

	for _, include := range rawConfig.Include {
		rawConfig := loadYamlFromPath(include)

		config.Processes = append(config.Processes, rawConfig.Processes...)
	}

	return config, nil
}

func loadYamlFromPath(path string) rawConfig {
	content := loadFileContent(ReplaceEnvVarsAndTilde(path))
	var rawConfig rawConfig
	if err := yaml.Unmarshal(content, &rawConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading YAML %s: %v\n", path, err)
		os.Exit(1)
	}

	return rawConfig
}

func loadFileContent(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load file path %s\n", path)
		os.Exit(1)
	}
	return data
}
