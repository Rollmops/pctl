package app

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"os"

	"gopkg.in/yaml.v2"
)

type YamlLoader struct {
}

type _rawConfig struct {
	Includes []string `yaml:"includes"`
	Config   `yaml:",inline"`
}

func NewYamlLoader() *YamlLoader {
	return &YamlLoader{}
}

func (l *YamlLoader) Load(path string) (*Config, error) {
	path = ExpandPath(path)
	rawConfig, err := loadYamlFromPath(path)
	if err != nil {
		return nil, err
	}

	err = loadIncludes(path, rawConfig.Includes, &rawConfig.Config)
	if err != nil {
		return nil, err
	}
	return &rawConfig.Config, nil
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
				return fmt.Errorf("config file %s is trying to include itself", baseConfigPath)
			}
			rawConfig, err := loadYamlFromPath(include)
			if err != nil {
				return err
			}
			config.ProcessConfigs = append(config.ProcessConfigs, rawConfig.ProcessConfigs...)
		}
	}
	return nil
}

func loadYamlFromPath(path string) (*_rawConfig, error) {
	content, err := loadFileContent(path)
	if err != nil {
		return nil, err
	}

	var rawConfig _rawConfig
	if err := yaml.UnmarshalStrict(content, &rawConfig); err != nil {
		return nil, fmt.Errorf("error reading YAML %s: %s", path, err.Error())
	}
	return &rawConfig, nil
}

func loadFileContent(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
