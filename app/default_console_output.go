package app

import (
	"fmt"
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

func (o *DefaultConsoleOutput) Write(processes []*Process) error {
	tw := table.NewWriter()
	// append a header row
	tw.AppendHeader(table.Row{"Name", "Status", "Pid", "Uptime", "Rss", "Vms", "Command"})

	// append some data rows
	runningCount := 0
	var rssSum uint64
	nowTime := time.Now()
	for _, p := range processes {
		pid := ""
		rss := ""
		vms := ""
		uptime := ""
		runningCommand := ""
		if p.IsRunning() {
			runningCount++
			uptime = _getUptime(p, nowTime)
			pid = strconv.Itoa(int(p.Info.GoPsutilProcess.Pid))
			memoryInfo, err := p.Info.GoPsutilProcess.MemoryInfo()
			if err != nil {
				rss = "error"
				vms = "error"
			} else {
				rssSum += memoryInfo.RSS
				rss = ByteCountIEC(memoryInfo.RSS)
				vms = ByteCountIEC(memoryInfo.VMS)
			}
			runningCommand = strings.Join(p.Info.RunningCommand, " ")
		}

		tw.AppendRow(table.Row{
			p.Config.Name,
			_getStatusString(p),
			pid,
			uptime,
			rss,
			vms,
			runningCommand,
		})
	}
	tw.AppendFooter(table.Row{
		"",
		fmt.Sprintf("Running: %d/%d", runningCount, len(processes)),
		"",
		"",
		fmt.Sprintf("Σ %s", ByteCountIEC(rssSum)),
	})
	tw.SetAutoIndex(true)
	tw.SetStyle(o.Style)
	tw.Style().Format.Header = text.FormatTitle
	tw.Style().Format.Footer = text.FormatTitle
	tw.Style().Options.SeparateColumns = false

	_, err := o.Writer.Write([]byte(tw.Render() + "\n"))
	return err
}

func _getUptime(p *Process, nowTime time.Time) string {
	var uptime string
	createTime, err := p.Info.GoPsutilProcess.CreateTime()
	if err != nil {
		uptime = "error"
	}
	uptimeInt := nowTime.Sub(time.Unix(createTime/1000, 0))
	uptime, err = DurationToString(uptimeInt)
	if err != nil {
		uptime = "error"
	}
	return uptime
}

func _getStatusString(p *Process) string {
	var statusString string
	if p.IsRunning() {
		statusString = OkColor("Running")
	} else {
		statusString = FailedColor("Stopped")
	}
	return statusString
}
