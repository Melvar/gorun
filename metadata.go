package main

import (
	"os"
	"gob"
)

type metadata struct {
	hash string
	lastused int64
	filename string
}

func readTable(filename string) (table map[string] metadata, err os.Error) {
	file, err := os.Open(filename, os.O_RDONLY | os.O_CREAT, perms)
	if err != nil {
		return
	}
	table = make(map[string] metadata)
	dec := gob.NewDecoder(file)
	for err == nil { //err is nil the first time, see above
		var entry metadata
		err = dec.Decode(&entry)
		if err == nil {
			table[entry.hash] = entry
		}
	}
	if err == os.EOF {
		err = nil
	}
	return
}

func writeTable(table map[string] metadata, filename string) (err os.Error) {
	file, err := os.Open(filename, os.O_WRONLY | os.O_CREAT, perms)
	if err != nil {
		return
	}
	enc := gob.NewEncoder(file)
	for _, entry := range table {
		err = enc.Encode(entry)
		if err != nil {
			return
		}
	}
	return
}
