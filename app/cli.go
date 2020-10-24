package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func CreateCliApp() (*cli.App, error) {
	filtersFlag := &cli.StringSliceFlag{
		Name:     "filter",
		Aliases:  []string{"f"},
		Required: false,
		Usage:    "IsRelevantForFilter output based on given label values. Format: <label>=<value>",
	}

	runningFlag := &cli.BoolFlag{
		Name:    "running",
		Aliases: []string{"r"},
		Value:   false,
		Usage:   "Filter only running processes (same as '--filter state.running==true')",
	}

	stoppedFlag := &cli.BoolFlag{
		Name:    "stopped",
		Aliases: []string{"s"},
		Value:   false,
		Usage:   "Filter only stopped processes (same as '--filter state.stopped==true')",
	}

	dirtyFlag := &cli.BoolFlag{
		Name:    "dirty",
		Aliases: []string{"d"},
		Value:   false,
		Usage:   "Filter only dirty processes (same as '--filter state.dirty==true')",
	}

	return &cli.App{
		Before: func(c *cli.Context) error {
			logLevelString := c.String("loglevel")
			level, err := logrus.ParseLevel(logLevelString)
			if err != nil {
				return fmt.Errorf("Unable to parse loglevel '%s'\n", logLevelString)
			}
			logrus.SetLevel(level)
			color.NoColor = c.Bool("no-color")
			CurrentContext.OutputWriter = os.Stdout
			return nil

		},
		Name:  "pctl",
		Usage: "process control",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "loglevel",
				Value:   "warning",
				EnvVars: []string{"PCTL_LOG_LEVEL"},
				Aliases: []string{"L"},
				Usage:   "level: trace,debug,info,warn,warning,error,fatal,panic",
			},
			&cli.BoolFlag{
				Name:    "no-color",
				Value:   false,
				EnvVars: []string{"PCTL_NO_COLOR"},
				Usage:   "do not use colors",
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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "comment",
						Usage: "add comment",
					},
					filtersFlag,
					runningFlag,
					stoppedFlag,
					dirtyFlag,
				},
				Action: func(c *cli.Context) error {
					filters := c.StringSlice("filter")
					filters = addShortcutFilters(c, filters)
					if c.NArg() == 0 && len(filters) == 0 {
						return fmt.Errorf("missing process names or filters")
					}
					return StartCommand(c.Args().Slice(), filters, c.String("comment"))
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
					filtersFlag,
					runningFlag,
					stoppedFlag,
					dirtyFlag,
				},
				Action: func(c *cli.Context) error {
					filters := c.StringSlice("filter")
					filters = addShortcutFilters(c, filters)
					if c.NArg() == 0 && len(filters) == 0 {
						return fmt.Errorf("missing process names or filters")
					}
					return StopCommand(c.Args().Slice(), filters, c.Bool("nowait"))
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
						Value:    "default",
						Usage: func() string {
							keys := make([]string, 0, len(FormatMap))
							for k := range FormatMap {
								keys = append(keys, k)
							}
							return "formats: " + strings.Join(keys, ",")
						}(),
					},
					filtersFlag,
					runningFlag,
					stoppedFlag,
					dirtyFlag,
				},
				Action: func(c *cli.Context) error {
					filters := c.StringSlice("filter")
					filters = addShortcutFilters(c, filters)
					format := c.String("format")
					if format == "" {
						format = "default"
					}
					return InfoCommand(c.Args().Slice(), format, filters)
				},
			},
		},
	}, nil
}

func addShortcutFilters(c *cli.Context, filters []string) []string {
	if c.Bool("running") {
		filters = append(filters, "state.running==true")
	}
	if c.Bool("stopped") {
		filters = append(filters, "state.stopped==true")
	}
	if c.Bool("dirty") {
		filters = append(filters, "state.dirty==true")
	}
	return filters
}
