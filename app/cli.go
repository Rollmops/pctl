package app

import (
	"fmt"
	"github.com/Rollmops/pctl/output"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func CreateCliApp(outputFile *os.File) *cli.App {
	return &cli.App{
		Before: func(c *cli.Context) error {
			logLevelString := c.String("loglevel")
			level, err := log.ParseLevel(logLevelString)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to parse loglevel '%s'\n", logLevelString)
				os.Exit(1)
			}
			log.SetLevel(level)

			CurrentContext, err = NewContext(outputFile)
			if err != nil {
				return err
			}
			err = CurrentContext.Initialize()
			if err != nil {
				return err
			}
			return CheckPersistenceConfigDiscrepancy()
		},
		Name:  "pctl",
		Usage: "process control",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "loglevel",
				Value:   "Info",
				EnvVars: []string{"PCTL_LOG_LEVEL"},
				Aliases: []string{"L"},
				Usage:   "level: trace,debug,info,warn,warning,error,fatal,panic",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "show the version information",
				Action: func(c *cli.Context) error {
					return fmt.Errorf("__TO_BE_IMPLEMENTED__")
				},
			},
			{
				Name:      "start",
				Usage:     "start process(es)",
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
				Usage:     "stop process(es)",
				ArgsUsage: "a list of process names",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "nowait",
						Value:   false,
						Usage:   "skip waiting until process stopped",
						EnvVars: []string{"PCTL_STOP_NO_WAIT"},
					},
					&cli.IntFlag{
						Name:    "wait-time",
						Value:   5,
						Usage:   "wait time in seconds",
						EnvVars: []string{"PCTL_STOP_WAIT_TIME"},
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("missing process names")
					}
					return StopCommand(c.Args().Slice(), c.Bool("nowait"), c.Int("wait-time"))
				},
			},
			{
				Name:      "info",
				Usage:     "show info for all or specified processes",
				ArgsUsage: "a list of process names - if empty, info of all processes will be shown",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "format",
						EnvVars:  []string{"PCTL_OUTPUT_FORMAT"},
						Required: false,
						Value:    "simple",
						Usage: func() string {
							keys := make([]string, 0, len(output.FormatMap))
							for k := range output.FormatMap {
								keys = append(keys, k)
							}
							return "formats: " + strings.Join(keys, ",")
						}(),
					},
				},
				Action: func(c *cli.Context) error {
					return InfoCommand(c.Args().Slice(), c.String("format"))
				},
			},
		},
	}
}
