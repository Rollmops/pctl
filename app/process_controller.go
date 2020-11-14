package app

type ProcessController interface {
	Start(names []string, filters Filters, comment string) error
	Stop(names []string, filters Filters, noWait bool, kill bool) error
	Restart(names []string, filters Filters, comment string, kill bool) error
	Kill(names []string, filters Filters) error
	Info(names []string, format string, filters Filters, columns []string) error
}
