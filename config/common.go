package config

const _configFileName string = "pctl.yml"

func (c *Config) FindByName(name string) *ProcessConfig {
	for _, p := range c.Processes {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

/*
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

	_, _ = fmt.Fprintf(os.Stderr, "Unable to to find valid config path: %v\n", possibleConfigPaths)
	os.Exit(1)

	return ""
}
*/
