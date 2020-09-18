package config

import (
	"fmt"
	"os"
	"path"
)

const _configFileName string = "pctl.yml"

func getConfigPath() string {
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
			return configPath
		}
	}

	fmt.Fprintf(os.Stderr, "Unable to to find valid config path: %v\n", possibleConfigPaths)
	os.Exit(1)

	return ""
}
