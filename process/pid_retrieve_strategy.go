package process

/*
Strategies:

- process is forking and running with a different Pid
  - Pid file
	- wait until Pid file was written
  - use cmdline/cmdline_suffix
    - wait until cmdline occurs in process list to retrieve new Pid
- in all other cases its sufficient to take the cmd.Process.Pid
*/

type PidRetrieveStrategy interface {
	Retrieve(*Process) (int32, error)
}

var PidRetrieveStrategies = map[string]PidRetrieveStrategy{}
