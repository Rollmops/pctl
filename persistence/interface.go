package persistence

type Data struct {
	Pid  int
	Name string
	Cmd  string
}

type Writer interface {
	Write(data []Data) error
}

type Reader interface {
	Read() ([]Data, error)
}
