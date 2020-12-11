package app

import (
	"github.com/fatih/color"
	"os"
)

var FormatMap = map[string]Output{}

var OkColor = color.New(color.FgGreen).SprintfFunc()
var FailedColor = color.New(color.FgRed).SprintfFunc()
var WarningColor = color.New(color.FgYellow).SprintfFunc()

type StatusReturn struct {
	OkMessage      string
	WarningMessage string
	Error          error
}

type Output interface {
	Write(list ProcessList, columnIds []string, sortColumns []string) error
	SetWriter(file *os.File)
}
