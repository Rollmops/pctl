package persistence

import (
	"encoding/csv"
	"os"
	"strconv"
)

type CsvWriter struct {
	path string
}

func NewCsvWriter(path string) CsvWriter {
	return CsvWriter{path: path}
}

func (c CsvWriter) Write(data []Data) error {
	file, err := os.OpenFile(c.path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, d := range data {
		err := writer.Write([]string{d.Name, strconv.Itoa(d.Pid), d.Cmd})
		if err != nil {
			return err
		}
	}
	return nil
}
