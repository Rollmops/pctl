package app

const (
	StateRunning = iota
	StateNotRunning
	StateStopped
	StateStarting
	StateStarted
	StateStopping
	StateStoppingError
	StateStartingError
	StateDependencyStoppingError
	StateKilling
	StateKilled
	StateKillingError
)

type ProcessStateChange struct {
	State   int
	Error   error
	Process *Process
}
