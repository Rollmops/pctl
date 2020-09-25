package config

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	"io/ioutil"
	"path/filepath"

	"os"

	"gopkg.in/yaml.v2"
)

func init() {
	SuffixConfigLoaderMap["yaml"] = &YamlLoader{}
	SuffixConfigLoaderMap["yml"] = &YamlLoader{}
}

type YamlLoader struct {
}

type _rawConfig struct {
	Includes     []string
	Processes    []*ProcessConfig
	StopStrategy StopStrategyConfig
}

func NewYamlLoader() *YamlLoader {
	return &YamlLoader{}
}

func (l *YamlLoader) Load(path string) (*Config, error) {
	path, err := common.ExpandPath(path)
	if err != nil {
		return nil, err
	}
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
	err = config.Validate()
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
				return fmt.Errorf("config file %s is trying to include itself", baseConfigPath)
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

func loadYamlFromPath(path string) (*_rawConfig, error) {
	content, err := loadFileContent(path)
	if err != nil {
		return nil, err
	}
	var rawConfig _rawConfig
	if err := yaml.Unmarshal(content, &rawConfig); err != nil {
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
