package persistence

import (
	"encoding/csv"
	"os"
	"strconv"
)

type CsvReader struct {
	path string
}

func NewCsvReader(path string) CsvReader {
	return CsvReader{path: path}
}

func (c CsvReader) Read() ([]Data, error) {
	file, err := os.OpenFile(c.path, os.O_RDONLY, 0755)
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
		pidStr, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		data = append(data, Data{Name: record[0], Pid: pidStr, Cmd: record[2]})
	}

	return data, nil

}
