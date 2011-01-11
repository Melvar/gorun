package main

import (
	"os"
	"io"
	"fmt"
	"log"
	hasht "crypto/sha256" //hash type, changeable
	db "./db" //primitive database adapter
)

func main() {

	if len(os.Args) < 2 { //no filename
		log.Exitln("No filename given")
	}
	scriptname := os.Args[1]
	scriptfile, err := os.Open(scriptname, os.O_RDONLY, 0)
	if err != nil {
		log.Exitln(err)
	}
	defer scriptfile.Close()

	hash := hasht.New()
	io.Copy(hash, scriptfile) //feed data to hash func
	hashstr := fmt.Sprintf("%x", hash.Sum()) //get hash as hex string

	table, err := db.UseFile(dbfilename) //get our data ready
	if err != nil {
		log.Exitln(err)
	}

	metadata, ok := table.Get(hashstr) //look for record of scriptfile
	compileAdd := func(add func(db.Table, db.Entry)) { //deduplication closure
		metadata, err = compile(scriptfile)
		if err != nil {
			log.Exitln(err)
		}
		add(table, metadata)
	}
	if !ok {
		compileAdd(db.Table.Add)
	} else {
		if !os.??fileexists(storedir + hashstr) {
			compileAdd(db.Table.Update)
		}
		metadata.lastused = os.Time()
		metadata.filename = scriptname
		table.Update(metadata)
	}
	if err := table.WriteBack(); err != nil {
		log.Println(err)
	}
//TODO: check for existence of executable
	os.Exec(hashstr, append(nil, scriptname, os.Args[2:]...), os.Environ())
}
