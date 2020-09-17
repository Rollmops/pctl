package app

import (
	"fmt"
	"os"
	"path"

	"github.com/urfave/cli/v2"
)

func CreateCliApp() *cli.App {
	return &cli.App{
		Name:  "pctl",
		Usage: "process control",
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "show the verion information",
				Action: func(c *cli.Context) error {
					fmt.Println("Version info")
					return nil
				},
			},
		},
	}
}

func GetConfigPath() string {
	cwd, _ := os.Getwd()
	configPaths := []string{
		path.Join(cwd, "pctl.yaml"),
		path.Join(os.Getenv("HOME"), "pctl.yaml"),
	}

	for _, configPath := range configPaths {
		if _, err := os.Stat(configPath); os.IsExist(err) {
			return configPath
		}
	}

	fmt.Fprintf(os.Stderr, "Unable to to find valid config path: %s\n", configPaths)
	os.Exit(1)

	return ""
}
