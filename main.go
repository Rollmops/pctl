package main

import (
	"fmt"
	"os"

	"github.com/Rollmops/pctl/app"
)

func main() {
	pctlApp := app.CreateCliApp()
	pctlApp.Run(os.Args)

	configPath := app.GetConfigPath()
	processMap := app.LoadProcessConfig(configPath)
	fmt.Println((*processMap)["b"].Cmd)
}
