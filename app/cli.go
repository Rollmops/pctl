package app

import (
	"fmt"
	"github.com/facebookgo/pidfile"
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
			noColor := c.Bool("no-color")
			color.NoColor = noColor
			logLevelString := c.String("loglevel")
			level, err := logrus.ParseLevel(logLevelString)
			if err != nil {
				return fmt.Errorf("Unable to parse loglevel %s\n", logLevelString)
			}
			logrus.SetLevel(level)
			logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: noColor})
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
				Name:  "agent",
				Usage: "pctl agent command group",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "pid-file",
						EnvVars: []string{"PCTL_AGENT_PID_FILE"},
						Usage:   "Path to the file to store and read the agent pid",
						Value:   "/tmp/pctl-agent.pid",
					},
				},
				Before: func(context *cli.Context) error {
					pidfile.SetPidfilePath(context.String("pid-file"))
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "start",
						Usage: "start the pctl agent",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "log-path",
								EnvVars: []string{"PCTL_AGENT_LOG_PATH"},
							},

							&cli.BoolFlag{
								Name:    "detach",
								Usage:   "detach mode - run the agent in the background",
								Aliases: []string{"d"},
								EnvVars: []string{"PCTL_AGENT_DETACH"},
							},
						},
						Action: func(c *cli.Context) error {
							return AgentStartCommand(c.String("log-path"), c.Bool("detach"))
						},
					},
					{
						Name:  "stop",
						Usage: "stop the pctl agent",
						Action: func(c *cli.Context) error {
							return AgentStopCommand()
						},
					},
					{
						Name:  "reload",
						Usage: "reload the configuration",
						Action: func(c *cli.Context) error {
							return AgentReloadCommand()
						},
					},
					{
						Name:  "status",
						Usage: "prints the agent status",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "derive-exit-code",
								Aliases: []string{"e"},
								Usage:   "set exit code to 1 if agent is not running",
							},
						},
						Action: func(c *cli.Context) error {
							return AgentStatusCommand(c.Bool("derive-exit-code"))
						},
					},
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
				Name:      "restart",
				Usage:     "restart process(es)",
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
				Action: func(c *cli.Context) error {
					filters := c.StringSlice("filter")
					filters = addShortcutFilters(c, filters)
					if c.NArg() == 0 && len(filters) == 0 {
						return fmt.Errorf("missing process names or filters")
					}
					return RestartCommand(c.Args().Slice(), filters, c.String("comment"))
				},
			},
			{
				Name:      "stop",
				Usage:     "stop process(es)",
				ArgsUsage: UsageNameSpecifiers,
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
				Name:      "kill",
				Usage:     "kill process(es)",
				ArgsUsage: UsageNameSpecifiers,
				Flags: []cli.Flag{
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
					columnsString := c.String("columns")
					columnsString = strings.ReplaceAll(columnsString, "+", "group,name,pid,status")
					columns := strings.Split(columnsString, ",")
					return InfoCommand(c.Args().Slice(), format, filters, columns)
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
