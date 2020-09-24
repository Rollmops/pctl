package app

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func CreateCliApp() *cli.App {
	return &cli.App{
		Name:  "pctl",
		Usage: "process control",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "loglevel",
				Value: "Info",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "show the version information",
				Action: func(c *cli.Context) error {
					fmt.Println("Version info")
					return nil
				},
			},
			{
				Name:      "start",
				Usage:     "start a process(es)",
				ArgsUsage: "a list of process names",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("missing process names")
					}
					return StartCommand(c.Args().Slice())
				},
			},
			{
				Name:      "stop",
				Usage:     "stop a process(es)",
				ArgsUsage: "a list of process names",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("missing process names")
					}
					return StopCommand(c.Args().Slice())
				},
			},
		},
		Before: func(c *cli.Context) error {
			logLevelString := c.String("loglevel")
			level, err := log.ParseLevel(logLevelString)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to parse loglevel '%s'\n", logLevelString)
				os.Exit(1)
			}
			log.SetLevel(level)

			CurrentContext, err = NewContext()
			if err != nil {
				return err
			}
			return CurrentContext.Initialize()
		},
	}
}
