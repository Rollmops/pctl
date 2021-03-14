package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

const __version__ = "__VERSION_PLACEHOLDER__"

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
		Usage:   "Filter only running processes (same as '--filter running==true')",
	}

	stoppedFlag := &cli.BoolFlag{
		Name:    "stopped",
		Aliases: []string{"s"},
		Value:   false,
		Usage:   "Filter only stopped processes (same as '--filter stopped==true')",
	}

	dirtyFlag := &cli.BoolFlag{
		Name:    "dirty",
		Aliases: []string{"d"},
		Value:   false,
		Usage:   "Filter only dirty processes (same as '--filter dirty==true')",
	}
	killFlag := &cli.BoolFlag{
		Name:    "kill",
		Value:   false,
		Usage:   "kill processes if unable to stop",
		EnvVars: []string{"PCTL_STOP_KILL"},
	}
	nowaitFlag := &cli.BoolFlag{
		Name:    "nowait",
		Value:   false,
		Usage:   "skip waiting until process stopped",
		EnvVars: []string{"PCTL_STOP_NO_WAIT"},
	}

	return &cli.App{
		Before: func(c *cli.Context) error {
			noColor := c.Bool("no-color")
			color.NoColor = noColor
			logLevelString := c.String("loglevel")
			level, err := logrus.ParseLevel(logLevelString)
			if err != nil {
				logrus.Fatalf("Unable to parse loglevel %s\n", logLevelString)
			}
			logrus.SetLevel(level)
			logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: noColor})
			CurrentContext.OutputWriter = os.Stdout
			CurrentContext.Context = c.String("context")
			err = CurrentContext.LoadConfig()
			if err != nil {
				logrus.Fatalf(err.Error())
			}
			return nil
		},
		Name:  "pctl",
		Usage: "process control",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "context",
				Value:   "default",
				EnvVars: []string{"PCTL_CONTEXT"},
				Aliases: []string{"c"},
				Usage:   "allows you to run independent pctl contexts in parallel",
			},
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
					fmt.Printf("pctl %s\n", __version__)
					return nil
				},
			},

			{
				Name:      "start",
				Usage:     "start process(es)",
				ArgsUsage: UsageNameSpecifiers,
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
				Before: func(_ *cli.Context) error {
					return CurrentContext.InitializeRunningProcessInfo()
				},
				Action: func(c *cli.Context) error {
					filters, err := createFilters(c)
					if err != nil {
						return err
					}
					if c.NArg() == 0 && len(filters) == 0 {
						return fmt.Errorf("missing process names or filters")
					}
					return StartCommand(c.Args().Slice(), filters, c.String("comment"))
				},
			},
			{
				Name:      "restart",
				Usage:     "restart process(es)",
				ArgsUsage: UsageNameSpecifiers,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "comment",
						Usage: "add comment",
					},
					killFlag,
					filtersFlag,
					runningFlag,
					stoppedFlag,
					dirtyFlag,
				},
				Before: func(_ *cli.Context) error {
					return CurrentContext.InitializeRunningProcessInfo()
				},
				Action: func(c *cli.Context) error {
					filters, err := createFilters(c)
					if err != nil {
						return err
					}
					if c.NArg() == 0 && len(filters) == 0 {
						return fmt.Errorf("missing process names or filters")
					}
					kill := c.Bool("kill")
					return RestartCommand(c.Args().Slice(), filters, c.String("comment"), kill)
				},
			},
			{
				Name:      "stop",
				Usage:     "stop process(es)",
				ArgsUsage: UsageNameSpecifiers,
				Flags: []cli.Flag{
					nowaitFlag,
					killFlag,
					filtersFlag,
					runningFlag,
					stoppedFlag,
					dirtyFlag,
				},
				Before: func(_ *cli.Context) error {
					return CurrentContext.InitializeRunningProcessInfo()
				},
				Action: func(c *cli.Context) error {
					filters, err := createFilters(c)
					if err != nil {
						return err
					}
					if c.NArg() == 0 && len(filters) == 0 {
						return fmt.Errorf("missing process names or filters")
					}
					noWait := c.Bool("nowait")
					kill := c.Bool("kill")
					if noWait && kill {
						return fmt.Errorf("unable to combine --nowait and --kill")
					}
					return StopCommand(c.Args().Slice(), filters, noWait, kill)
				},
			},
			{
				Name:      "kill",
				Usage:     "kill process(es)",
				ArgsUsage: UsageNameSpecifiers,
				Flags: []cli.Flag{
					filtersFlag,
					runningFlag,
					stoppedFlag,
					dirtyFlag,
				},
				Before: func(_ *cli.Context) error {
					return CurrentContext.InitializeRunningProcessInfo()
				},
				Action: func(c *cli.Context) error {
					filters, err := createFilters(c)
					if err != nil {
						return err
					}
					if c.NArg() == 0 && len(filters) == 0 {
						return fmt.Errorf("missing process names or filters")
					}
					return KillCommand(c.Args().Slice(), filters)
				},
			},
			{
				Name:      "info",
				Usage:     "show info for all or specified processes",
				ArgsUsage: "a list of process names - if empty, info of all processes will be shown",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "format",
						EnvVars:  []string{"PCTL_INFO_FORMAT"},
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
					&cli.StringFlag{
						Name:     "columns",
						EnvVars:  []string{"PCTL_INFO_COLUMNS"},
						Required: false,
						Aliases:  []string{"c"},
						Value:    "group,name,pid,status",
						Usage: func() string {
							columns := make([]string, 0, len(PropertyMap))
							for k := range PropertyMap {
								columns = append(columns, k)
							}
							return "comma separated list of columns, available: " + strings.Join(columns, ",")
						}(),
					},
					&cli.StringFlag{
						Name:     "sort",
						EnvVars:  []string{"PCTL_INFO_SORT_COLUMNS"},
						Required: false,
						Usage:    "comma separated list of the columns to sort",
					},
					filtersFlag,
					runningFlag,
					stoppedFlag,
					dirtyFlag,
				},
				Before: func(_ *cli.Context) error {
					return CurrentContext.InitializeRunningProcessInfo()
				},
				Action: func(c *cli.Context) error {
					filters, err := createFilters(c)
					if err != nil {
						return err
					}
					format := c.String("format")
					if format == "" {
						format = "default"
					}
					columnsString := c.String("columns")
					columnsString = strings.ReplaceAll(columnsString, "+", "group,name,pid,status")
					columns := strings.Split(columnsString, ",")
					sortColumns := strings.Split(c.String("sort"), ",")
					return InfoCommand(c.Args().Slice(), format, filters, columns, sortColumns)
				},
			},
		},
	}, nil
}

func createFilters(c *cli.Context) (Filters, error) {
	filters, err := NewFilters(c.StringSlice("filter"))
	if err != nil {
		return nil, err
	}
	filters = addShortcutFilters(c, filters)
	return filters, nil
}

func addShortcutFilters(c *cli.Context, filters Filters) Filters {
	if c.Bool("running") {
		filters = append(filters, &Filter{
			field:    "running",
			operator: "=",
			value:    "true",
		})
	}
	if c.Bool("stopped") {
		filters = append(filters, &Filter{
			field:    "stopped",
			operator: "=",
			value:    "true",
		})
	}
	if c.Bool("dirty") {
		filters = append(filters, &Filter{
			field:    "dirty",
			operator: "=",
			value:    "true",
		})
	}
	return filters
}
