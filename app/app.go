package app

func Run(args []string) error {
	pctlApp, err := CreateCliApp()
	if err != nil {
		return err
	}

	err = CurrentContext.Initialize()
	if err != nil {
		return err
	}
	err = ValidateAcyclicDependencies()
	if err != nil {
		return err
	}

	return pctlApp.Run(args)
}
