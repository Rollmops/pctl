package main

import (
	"github.com/Rollmops/pctl/app"
	"log"
	"os"
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
