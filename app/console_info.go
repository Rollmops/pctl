package app

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"io"
	"os"
	"strings"
)

func init() {
	FormatMap["default"] = &DefaultConsoleOutput{Style: table.StyleBold}
	FormatMap["bright"] = &DefaultConsoleOutput{Style: table.StyleColoredBright}
	FormatMap["simple"] = &DefaultConsoleOutput{Style: table.StyleDefault}
	FormatMap["dark"] = &DefaultConsoleOutput{Style: table.StyleColoredDark}
	FormatMap["black-on-green-white"] = &DefaultConsoleOutput{Style: table.StyleColoredBlackOnGreenWhite}
	FormatMap["black-on-blue-white"] = &DefaultConsoleOutput{Style: table.StyleColoredBlackOnBlueWhite}
	FormatMap["black-on-cyan-white"] = &DefaultConsoleOutput{Style: table.StyleColoredBlackOnCyanWhite}
	FormatMap["black-on-magenta-white"] = &DefaultConsoleOutput{Style: table.StyleColoredBlackOnMagentaWhite}
	FormatMap["black-on-red-white"] = &DefaultConsoleOutput{Style: table.StyleColoredBlackOnRedWhite}
	FormatMap["black-on-yellow-white"] = &DefaultConsoleOutput{Style: table.StyleColoredBlackOnYellowWhite}
	FormatMap["blue-white-on-black"] = &DefaultConsoleOutput{Style: table.StyleColoredBlueWhiteOnBlack}
	FormatMap["bright"] = &DefaultConsoleOutput{Style: table.StyleColoredBright}
	FormatMap["green-white-on-black"] = &DefaultConsoleOutput{Style: table.StyleColoredGreenWhiteOnBlack}
	FormatMap["red-white-on-black"] = &DefaultConsoleOutput{Style: table.StyleColoredRedWhiteOnBlack}
	FormatMap["cyan-white-on-black"] = &DefaultConsoleOutput{Style: table.StyleColoredCyanWhiteOnBlack}
	FormatMap["yellow-white-on-black"] = &DefaultConsoleOutput{Style: table.StyleColoredYellowWhiteOnBlack}
	FormatMap["rounded"] = &DefaultConsoleOutput{Style: table.StyleRounded}
}

type DefaultConsoleOutput struct {
	Writer io.Writer
	Style  table.Style
}

func (o *DefaultConsoleOutput) SetWriter(writer *os.File) {
	o.Writer = writer
}

func (o *DefaultConsoleOutput) Write(processes ProcessList, columnIds []string, sortColumns []string) error {
	tw := table.NewWriter()
	tw.SetStyle(o.Style)
	tw.Style().Format.Header = text.FormatTitle
	tw.Style().Format.Footer = text.FormatTitle
	tw.Style().Options.SeparateColumns = false
	tw.SetAutoIndex(true)

	var header table.Row
	var footer table.Row
	for _, columnId := range columnIds {
		property := PropertyMap[columnId]
		if property != nil {
			header = append(header, property.Name())
			footerValue, err := property.FormattedSumValue(processes)
			if err != nil {
				return err
			}
			footer = append(footer, footerValue)
		} else {
			return fmt.Errorf("column %s not available", columnId)
		}
	}
	tw.AppendHeader(header)
	tw.AppendFooter(footer)

	for _, process := range processes {
		var row table.Row
		for _, columnId := range columnIds {
			property := PropertyMap[columnId]
			if property != nil {
				value, err := property.Value(process, true)
				if err != nil {
					return err
				}
				row = append(row, value)
			}
		}
		tw.AppendRow(row)
	}
	tw.SortBy(_createSortBySlice(sortColumns, columnIds))
	_, err := o.Writer.Write([]byte(tw.Render() + "\n"))
	return err
}

func _createSortBySlice(sortColumns []string, columnIds []string) []table.SortBy {
	var sortBySlice []table.SortBy
	for _, sortColumn := range sortColumns {
		for i, columnId := range columnIds {
			if strings.ToLower(columnId) == strings.ToLower(sortColumn) {
				sortBySlice = append(sortBySlice, table.SortBy{Number: i + 1})
			}
		}
	}
	return sortBySlice
}
