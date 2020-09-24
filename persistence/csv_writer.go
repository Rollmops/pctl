package persistence

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

type CsvWriter struct {
	stateFilePath string
}

func NewCsvReader() (*CsvReader, error) {
	stateFilePath, err := GetStateFilePath()
	if err != nil {
		return nil, err
	}
	stateFilePath = filepath.Join(stateFilePath, "state.csv")
	return &CsvReader{stateFilePath: stateFilePath}, nil
}

func NewTestCsvReader(stateFilePath string) CsvReader {
	return CsvReader{stateFilePath: stateFilePath}
}

func (c *CsvWriter) Write(data []Data) error {
	file, err := os.OpenFile(c.stateFilePath, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, d := range data {
		err := writer.Write([]string{d.Name, strconv.Itoa(int(d.Pid)), d.Cmd})
		if err != nil {
			return err
		}
	}
	return nil
}
