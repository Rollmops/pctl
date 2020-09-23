package main

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
)

func main() {
	//pctlApp := app.CreateCliApp()
	//_ = pctlApp.Run(os.Args)
	//log.Debug("Starting application pctl")

	c := config.ProcessConfig{
		Name:                    "test",
		Cmd:                     []string{"test/fixtures/daemon_test.py"},
		PidRetrieveStrategyName: "cmdline",
	}

	p := process.NewProcess(c)
	_ = p.Start()

	fmt.Println(p.Pid())

}
