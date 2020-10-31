package app

type ScriptStopStrategyConfig struct {
	Script `yaml:",inline"`
}

type ScriptStopStrategy struct {
	ScriptStopStrategyConfig
}

func (s *ScriptStopStrategy) Stop(process *Process) error {
	cmd, err := s.CreateCommand(process)
	if err != nil {
		return err
	}
	return cmd.Run()
}
