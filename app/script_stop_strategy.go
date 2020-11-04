package app

type ScriptStopStrategy struct {
	Exec `yaml:",inline"`
}

func (s *ScriptStopStrategy) Stop(process *Process) error {
	cmd, err := s.CreateCommand(process)
	if err != nil {
		return err
	}
	return cmd.Run()
}
