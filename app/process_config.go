package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

type ProcessConfigType struct {
	Cmd string `yaml:"cmd"`
}

type ProcessConfigTypeMap map[string]*ProcessConfigType

var processConfigMap ProcessConfigTypeMap

func LoadProcessConfig(configPath string) *ProcessConfigTypeMap {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load file path %s\n", configPath)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(data, &processConfigMap); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading YAML %s: %v\n", configPath, err)
		os.Exit(1)
	}
	return &processConfigMap
}
