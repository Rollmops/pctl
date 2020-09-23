package config

type Loader interface {
	Load(path string) (*Config, error)
}

type ProcessConfig struct {
	Name                    string
	Cmd                     []string
	PidRetrieveStrategyName string
}

type Config struct {
	Processes []ProcessConfig
}
