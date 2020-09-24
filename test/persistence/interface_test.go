package persistence_test

import (
	"github.com/Rollmops/pctl/persistence"
	"testing"
)

func TestAddOrUpdateEntry(t *testing.T) {

	data := persistence.Data{
		Entries: []persistence.DataEntry{
			{
				Pid:  100,
				Name: "p1",
				Cmd:  "sleep 100",
			},
		},
	}

	data.AddOrUpdateEntry(&persistence.DataEntry{
		Name: "p2",
		Pid:  101,
		Cmd:  "ls -la",
	})

	if len(data.Entries) != 2 {
		t.Fatalf("2 != %d", len(data.Entries))
	}

	data.AddOrUpdateEntry(&persistence.DataEntry{
		Name: "p2",
		Pid:  102,
		Cmd:  "ls -la",
	})

	if len(data.Entries) != 2 {
		t.Fatalf("2 != %d", len(data.Entries))
	}

	p2 := data.FindByName("p2")
	if p2 == nil {
		t.Fatal("could not find p2")
	}

	if p2.Pid != 102 {
		t.Fatal("pid was not updated")
	}

}
