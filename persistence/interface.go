package persistence

type Data struct {
	Entries []DataEntry
}

type DataEntry struct {
	Pid  int32
	Name string
	Cmd  string
}

type Writer interface {
	Write(data *Data) error
}

type Reader interface {
	Read() (*Data, error)
}

func (d *Data) FindByName(name string) *DataEntry {
	for _, entry := range d.Entries {
		if entry.Name == name {
			return &entry
		}
	}
	return nil
}

func (d *Data) AddOrUpdateEntry(entry *DataEntry) {
	for index, _entry := range d.Entries {
		if _entry.Name == entry.Name {
			d.Entries[index].Pid = entry.Pid
			d.Entries[index].Cmd = entry.Cmd
			return
		}
	}
	d.Entries = append(d.Entries, *entry)
}
