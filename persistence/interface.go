package persistence

type Data struct {
	Pid  int
	Name string
	Cmd  string
}

type Writer interface {
	write(data []Data) error
}

type Reader interface {
	read() ([]Data, error)
}
