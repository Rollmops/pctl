package process

/*
Strategies:

- process is forking and running with a different pid
  - pid file
	- wait until pid file was written
  - use cmdline/cmdline_suffix
    - wait until cmdline occurs in process list to retrieve new pid
- in all other cases its sufficient to take the cmd.Process.Pid
*/

type PidRetrieveStrategy interface {
	Retrieve(*Process) (int32, error)
}

var PidRetrieveStrategies = map[string]PidRetrieveStrategy{}
