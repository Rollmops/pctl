package persistence_test

import (
	"github.com/Rollmops/pctl/persistence"
	"testing"
)

func TestDataHandling(t *testing.T) {

	data := persistence.Data{
		Entries: []*persistence.DataEntry{
			{
				Pid:     100,
				Name:    "p1",
				Command: []string{"sleep", "100"},
			},
		},
	}

	data.AddOrUpdateEntry(&persistence.DataEntry{
		Name:    "p2",
		Pid:     101,
		Command: []string{"ls", "-la"},
	})

	if len(data.Entries) != 2 {
		t.Fatalf("2 != %d", len(data.Entries))
	}

	data.AddOrUpdateEntry(&persistence.DataEntry{
		Name:    "p2",
		Pid:     102,
		Command: []string{"ls", "-la"},
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

	data.RemoveByName("p1")
	if len(data.Entries) != 1 {
		t.Fatalf("1 != %d", len(data.Entries))
	}
	if data.Entries[0].Name != "p2" {
		t.Fatalf("p2 != %s", data.Entries[0].Name)
	}
}
