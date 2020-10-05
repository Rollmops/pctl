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

type Writer struct {
	stateFilePath string
}

func NewCsvReader() (*Reader, error) {
	stateFilePath, err := persistence.GetStateFilePath()
	if err != nil {
		return nil, err
	}
	stateFilePath = filepath.Join(stateFilePath, "state.csv")
	return &Reader{stateFilePath: stateFilePath}, nil
}

func NewTestCsvReader(stateFilePath string) Reader {
	return Reader{stateFilePath: stateFilePath}
}

func (c *Writer) Write(data *persistence.Data) error {
	log.Debugf("Opening %s", c.stateFilePath)
	file, err := os.OpenFile(c.stateFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	log.Debugf("Writing %d entries to %s", len(data.Entries), c.stateFilePath)
	for _, d := range data.Entries {
		commandString, err := json.Marshal(d.Command)
		if err != nil {
			return err
		}
		err = writer.Write([]string{
			d.Name,
			strconv.Itoa(int(d.Pid)),
			string(commandString),
			d.Comment,
			strconv.Itoa(d.MarkFlag),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
