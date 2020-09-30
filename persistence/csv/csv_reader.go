package csv

import (
	"encoding/csv"
	"encoding/json"
	"github.com/Rollmops/pctl/persistence"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
)

type CsvReader struct {
	stateFilePath string
}

func NewCsvWriter() (*CsvWriter, error) {
	stateFilePath, err := persistence.GetStateFilePath()
	if err != nil {
		return nil, err
	}
	stateFilePath = filepath.Join(stateFilePath, "state.csv")
	log.Tracef("Creating csv writer with state file path: %s", stateFilePath)
	return &CsvWriter{stateFilePath: stateFilePath}, nil
}

func NewTestCsvWriter(stateFilePath string) CsvWriter {
	return CsvWriter{stateFilePath: stateFilePath}
}

func (c *CsvReader) Read() (*persistence.Data, error) {
	log.Debugf("Opening %s", c.stateFilePath)
	file, err := os.OpenFile(c.stateFilePath, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	defer func() {
		log.Tracef("Closing %s", c.stateFilePath)
		_ = file.Close()
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []*persistence.DataEntry
	log.Debugf("Reading %d records from state file", len(records))
	for _, record := range records {
		data, err = _readAndAppendCsvRecord(data, record)
		if err != nil {
			return nil, err
		}
	}
	return &persistence.Data{Entries: data}, nil
}

func _readAndAppendCsvRecord(data []*persistence.DataEntry, record []string) ([]*persistence.DataEntry, error) {
	log.Tracef("Reading csv record %v", record)
	pid, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, err
	}
	var command []string
	err = json.Unmarshal([]byte(record[2]), &command)
	if err != nil {
		return nil, err
	}
	data = append(data, &persistence.DataEntry{Name: record[0], Pid: int32(pid), Command: command, Comment: record[3]})
	return data, nil
}
