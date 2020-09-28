package output

import (
	"fmt"
	"github.com/fatih/color"
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

func PrintMessageAndStatus(message string, function func() StatusReturn) error {
	fmt.Printf("%s ... ", message)
	statusReturn := function()
	if statusReturn.Error != nil {
		fmt.Printf("%s\n", FailedColor("Failed (%s)", statusReturn.Error.Error()))
		return statusReturn.Error
	}
	if statusReturn.WarningMessage != "" {
		fmt.Printf("%s\n", WarningColor("Warning (%s)", statusReturn.WarningMessage))
		return nil
	}
	okMessageString := ""
	if statusReturn.OkMessage != "" {
		okMessageString = fmt.Sprintf(" (%s)", statusReturn.OkMessage)
	}
	fmt.Printf("%s\n", OkColor("Ok%s", okMessageString))
	return nil
}
