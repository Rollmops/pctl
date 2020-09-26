package persistence

type Data struct {
	Entries []*DataEntry
}

type DataEntry struct {
	Pid     int32
	Name    string
	Command []string
}

type Writer interface {
	Write(data *Data) error
}

type Reader interface {
	Read() (*Data, error)
}

// Find the persistence data entry for that name or nil
func (d *Data) FindByName(name string) *DataEntry {
	for _, entry := range d.Entries {
		if entry.Name == name {
			return entry
		}
	}
	return nil
}

func (d *Data) AddOrUpdateEntry(entry *DataEntry) {
	for index, _entry := range d.Entries {
		if _entry.Name == entry.Name {
			d.Entries[index].Pid = entry.Pid
			d.Entries[index].Command = entry.Command
			return
		}
	}
	d.Entries = append(d.Entries, entry)
}

func (d *Data) RemoveByName(name string) {
	index := d._findIndex(name)
	if index != -1 {
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
