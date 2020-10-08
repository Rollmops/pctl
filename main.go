package main

import (
	"github.com/Rollmops/pctl/app"
	"log"
	"os"
)

func main() {
	pctlApp, err := app.CreateCliApp()
	if err != nil {
		log.Fatal(err)
	}
	err = pctlApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
