package app

func RestartCommand(names []string, filters []string, comment string) error {
	err := StopCommand(names, filters, false)
	if err != nil {
		return err
	}
	return StartCommand(names, filters, comment)
}
