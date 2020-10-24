package output

import (
	"os"
)

type Output interface {
	Write([]*Info) error
	SetWriter(file *os.File)
}
