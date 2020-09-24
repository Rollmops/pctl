package persistence_test

import (
	"github.com/Rollmops/pctl/persistence"
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
	csvWriter := persistence.NewTestCsvWriter(tmpFile.Name())

	data := []persistence.Data{
		{
			Pid:  1,
			Name: "process1",
			Cmd:  "sleep infinity",
		},
		{
			Pid:  2,
			Name: "process2",
			Cmd:  "cat",
		},
	}

	err = csvWriter.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	csvReader := persistence.NewTestCsvReader(tmpFile.Name())

	readData, err := csvReader.Read()
	if err != nil {
		t.Fatal(err)
	}

	for index := range data {
		if data[index] != readData[index] {
			t.Fatalf("%v != %v", data[index], readData[index])
		}
	}

}
