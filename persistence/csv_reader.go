package persistence

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

type CsvReader struct {
	stateFilePath string
}

func NewCsvWriter() (*CsvWriter, error) {
	stateFilePath, err := GetStateFilePath()
	if err != nil {
		return nil, err
	}
	stateFilePath = filepath.Join(stateFilePath, "state.csv")
	return &CsvWriter{stateFilePath: stateFilePath}, nil
}

func NewTestCsvWriter(stateFilePath string) CsvWriter {
	return CsvWriter{stateFilePath: stateFilePath}
}

func (c *CsvReader) Read() ([]Data, error) {
	file, err := os.OpenFile(c.stateFilePath, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []Data
	for _, record := range records {
		pid, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		data = append(data, Data{Name: record[0], Pid: int32(pid), Cmd: record[2]})
	}

	return data, nil

}
