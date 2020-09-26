package output

import (
	"fmt"
	"github.com/fatih/color"
)

var FormatMap = map[string]Output{}

var Green = color.New(color.FgGreen).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()

func PrintMessageAndStatus(message string, function func() error) error {
	fmt.Printf("%s ... ", message)
	err := function()
	if err != nil {
		fmt.Printf("%s\n", Red("Failed"))
		return err
	}
	fmt.Printf("%s\n", Green("Ok"))
	return nil
}
