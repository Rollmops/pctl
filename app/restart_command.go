package app

import "fmt"

func RestartCommand(names []string, filters []string, comment string, kill bool) error {
	err := StopCommand(names, filters, false, kill)
	if err != nil {
		return err
	}
	err = CurrentContext.Cache.Refresh()
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	if err != nil {
		return err
	}

	// TODO remove filters that may change during restart (e.g. dirty)
	return StartCommand(names, filters, comment)
}
