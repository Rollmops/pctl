package csv_test

import (
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/persistence"
	"github.com/Rollmops/pctl/persistence/csv"
	"io/ioutil"
	"os"
	"testing"
)

func TestWriteReadCsv(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "pctl_test.*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	csvWriter := csv.NewTestCsvWriter(tmpFile.Name())

	data := &persistence.Data{
		Entries: []*persistence.DataEntry{
			{
				Pid:     1,
				Name:    "process1",
				Command: []string{"sleep", "infinity"},
			},
			{
				Pid:     2,
				Name:    "process2",
				Command: []string{"cat"},
			},
		},
	}

	err = csvWriter.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	csvReader := csv.NewTestCsvReader(tmpFile.Name())

	readData, err := csvReader.Read()
	if err != nil {
		t.Fatal(err)
	}

	for index := range data.Entries {
		if data.Entries[index].Pid != readData.Entries[index].Pid ||
			!common.CompareStringSlices(data.Entries[index].Command, readData.Entries[index].Command) ||
			data.Entries[index].Name != data.Entries[index].Name {
			t.Fatalf("%v != %v", data.Entries[index], readData.Entries[index])
		}
	}

}
