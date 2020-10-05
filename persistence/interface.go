package persistence

import log "github.com/sirupsen/logrus"

type Data struct {
	Entries []*DataEntry
}

const (
	MarkedAsStopped = iota
	MarkedAsStarted
)

type DataEntry struct {
	Pid      int32
	Name     string
	Command  []string
	Comment  string
	MarkFlag int
}

type Writer interface {
	Write(data *Data) error
}

type Reader interface {
	Read() (*Data, error)
}

// Find the persistence data entry for that name or nil
func (d *Data) FindByName(name string) *DataEntry {
	log.Tracef("Finding persistence data entry for '%s'", name)
	for _, entry := range d.Entries {
		if entry.Name == name {
			log.Tracef("Found persistence data entry %v", entry)
			return entry
		}
	}
	log.Tracef("Did not find a persistence data entry for '%s'", name)
	return nil
}

func (d *Data) AddOrUpdateEntry(entry *DataEntry) {
	for index, _entry := range d.Entries {
		if _entry.Name == entry.Name {
			log.Tracef("Updating persistence data entry %v", entry)
			d.Entries[index].Pid = entry.Pid
			d.Entries[index].Command = entry.Command
			d.Entries[index].Comment = entry.Comment
			d.Entries[index].MarkFlag = entry.MarkFlag
			return
		}
	}
	log.Tracef("Adding new persistence data entry %v", entry)
	d.Entries = append(d.Entries, entry)
}

func (d *Data) RemoveByName(name string) {
	index := d._findIndex(name)
	if index != -1 {
		log.Tracef("Removing persistence data entry %v", d.Entries[index])
		d.Entries[index] = d.Entries[len(d.Entries)-1]
		d.Entries = d.Entries[:len(d.Entries)-1]
	}
}

func (d *Data) _findIndex(name string) int {
	for index, entry := range d.Entries {
		if entry.Name == name {
			return index
		}
	}
	return -1
}
