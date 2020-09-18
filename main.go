package main

import (
	"os"

	"github.com/Rollmops/pctl/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	pctlApp := app.CreateCliApp()
	_ = pctlApp.Run(os.Args)
	log.Debug("Starting application pctl")

}
