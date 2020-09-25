package config

type Loader interface {
	Load(path string) (*Config, error)
}

type ProcessConfig struct {
	Name                    string
	Command                 []string `yaml:"cmd"`
	PidRetrieveStrategyName string   `yaml:"pidStrategy"`
}

type Config struct {
	Processes []*ProcessConfig
}
