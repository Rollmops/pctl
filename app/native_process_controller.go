package app

import (
	"fmt"
	"github.com/Songmu/prompter"
)

type NativeProcessController struct{}

func (p *NativeProcessController) Start(names []string, filters Filters, comment string) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}

	if CurrentContext.Config.PromptForStart && !prompter.YN(fmt.Sprintf("Do you really want to start?"), false) {
		return nil
	}

	return StartProcesses(processes, comment)
}

func (p *NativeProcessController) Stop(names []string, filters Filters, noWait bool, kill bool) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}

	if CurrentContext.Config.PromptForStop && !prompter.YN(fmt.Sprintf("Do you really want to proceed stopping?"), false) {
		return nil
	}
	_, err = StopProcesses(processes, noWait, kill)
	return err
}
func (p *NativeProcessController) Restart(names []string, filters Filters, comment string, kill bool) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}
	if CurrentContext.Config.PromptForStop && !prompter.YN(fmt.Sprintf("Do you really want to proceed with restart?"), false) {
		return nil
	}
	stoppedProcesses, err := StopProcesses(processes, false, kill)
	if err != nil {
		return err
	}
	if len(stoppedProcesses) == 0 {
		return nil
	}
	err = CurrentContext.Cache.Refresh()
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	if err != nil {
		return err
	}

	return StartProcesses(stoppedProcesses, comment)
}
func (p *NativeProcessController) Kill(names []string, filters Filters) error {
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, len(filters) > 0)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}

	if !prompter.YN(fmt.Sprintf("Do you really want to proceed killing?"), false) {
		return nil
	}

	return KillProcesses(processes)
}
func (p *NativeProcessController) Info(names []string, format string, filters Filters, columns []string) error {
	o := FormatMap[format]
	if o == nil {
		return fmt.Errorf("unknown format: %s", format)
	}
	o.SetWriter(CurrentContext.OutputWriter)
	processes, err := CurrentContext.Config.CollectProcessesByNameSpecifiers(names, filters, true)
	if err != nil {
		return err
	}
	if len(processes) == 0 {
		return fmt.Errorf(MsgNoMatchingProcess)
	}
	return o.Write(processes, columns)
}
