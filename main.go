package main

import (
	"github.com/Rollmops/pctl/app"
	"log"
	"os"
)

func main() {
	pctlApp := app.CreateCliApp(os.Stdout)
	err := pctlApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
