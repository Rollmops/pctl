package main

import (
	"fmt"
	"os"

	"github.com/Rollmops/pctl/app"
	"github.com/Rollmops/pctl/config"

	log "github.com/sirupsen/logrus"
)

func main() {
	pctlApp := app.CreateCliApp()
	_ = pctlApp.Run(os.Args)
	log.Debug("Starting application pctl")

	_config, _ := config.LoadConfig("")

	fmt.Println(_config.Processes)
}
