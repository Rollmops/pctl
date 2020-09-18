package persistence

import (
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
	csvWriter := NewCsvWriter(tmpFile.Name())

	data := []Data{
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

	err = csvWriter.write(data)
	if err != nil {
		t.Fatal(err)
	}

	csvReader := NewCsvReader(tmpFile.Name())

	readData, err := csvReader.read()
	if err != nil {
		t.Fatal(err)
	}

	for index, _ := range data {
		if data[index] != readData[index] {
			t.Fatalf("%v != %v", data[index], readData[index])
		}
	}

}
