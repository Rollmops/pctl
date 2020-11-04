package main

import (
	"github.com/Rollmops/pctl/app"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
