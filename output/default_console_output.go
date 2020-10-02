package output

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	FormatMap["default"] = &DefaultConsoleOutput{Style: table.StyleColoredBright}
	FormatMap["simple"] = &DefaultConsoleOutput{Style: table.StyleDefault}
	FormatMap["bold"] = &DefaultConsoleOutput{Style: table.StyleBold}
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

func (o *DefaultConsoleOutput) Write(infoEntries []*InfoEntry) error {
	tw := table.NewWriter()
	// append a header row
	tw.AppendHeader(table.Row{"Name", "Status", "Pid", "Uptime", "Rss", "Vms", "Command"})

	// append some data rows
	runningCount := 0
	var rssSum uint64
	nowTime := time.Now()
	for _, entry := range infoEntries {
		if entry.IsRunning {
			runningCount++
		}
		uptime := _getUptime(entry, nowTime)
		pid := ""
		rss := ""
		vms := ""
		if entry.RunningInfo != nil {
			pid = strconv.Itoa(int(entry.RunningInfo.Pid))
			memoryInfo, err := entry.RunningInfo.MemoryInfo()
			if err != nil {
				rss = "error"
				vms = "error"
			} else {
				rssSum += memoryInfo.RSS
				rss = common.ByteCountIEC(memoryInfo.RSS)
				vms = common.ByteCountIEC(memoryInfo.VMS)
			}
		}

		tw.AppendRow(table.Row{
			entry.Name,
			_getStatusString(entry),
			pid,
			uptime,
			rss,
			vms,
			strings.Join(entry.RunningCommand, " "),
		})
	}
	tw.AppendFooter(table.Row{
		"",
		fmt.Sprintf("Running: %d/%d", runningCount, len(infoEntries)),
		"",
		"",
		fmt.Sprintf("Σ %s", common.ByteCountIEC(rssSum)),
	})
	tw.SetAutoIndex(true)
	// sort by last name and then by salary
	//tw.SortBy([]table.SortBy{{Name: "", Mode: table.Dsc}, {Name: "Salary", Mode: table.AscNumeric}})
	// use a ready-to-use Style
	tw.SetStyle(o.Style)
	// customize the Style and change some stuff
	tw.Style().Format.Header = text.FormatTitle
	//tw.Style().Format.Row = text.FormatLower
	tw.Style().Format.Footer = text.FormatTitle
	tw.Style().Options.SeparateColumns = false

	_, err := o.Writer.Write([]byte(tw.Render() + "\n"))
	return err
}

func _getUptime(entry *InfoEntry, nowTime time.Time) string {
	var uptime string
	if entry.RunningInfo != nil {
		createTime, err := entry.RunningInfo.CreateTime()
		if err != nil {
			uptime = "error"
		}
		uptimeInt := nowTime.Sub(time.Unix(createTime/1000, 0))
		uptime, err = common.DurationToString(uptimeInt)
		if err != nil {
			uptime = "error"
		}
	}
	return uptime
}

func _getStatusString(entry *InfoEntry) string {
	var statusString string
	if entry.IsRunning {
		statusString = OkColor("Running")
	} else if entry.StoppedUnexpectedly {
		statusString = FailedColor("Crashed")
	} else {
		statusString = WarningColor("Stopped")
	}
	return statusString
}
