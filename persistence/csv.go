package persistence

import (
	"encoding/csv"
	"os"
	"strconv"
)

type CsvWriter struct {
	path string
}

type CsvReader struct {
	path string
}

func NewCsvWriter(path string) CsvWriter {
	return CsvWriter{path: path}
}

func NewCsvReader(path string) CsvReader {
	return CsvReader{path: path}
}

func (c CsvWriter) write(data []Data) error {
	file, err := os.OpenFile(c.path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

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

func (c CsvReader) read() ([]Data, error) {
	file, err := os.OpenFile(c.path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

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
