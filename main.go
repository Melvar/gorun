package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"path"
	hasht "crypto/sha256" //hash type, changeable
)

const (
	dbfilename = path.Join(os.Getenv("HOME"), ".gorundb.gob")
	storedir = path.Join(os.Getenv("HOME"), ".gorun")
	dirperms = 0755
)

func main() {

	err := os.MkdirAll(storedir, dirperms)
	if err != nil {
		log.Fatalln(err)
	}

	if len(os.Args) < 2 { //no filename
		log.Fatalln("No filename given")
	}
	scriptname := os.Args[1]
	scriptfile, err := os.Open(scriptname, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalln(err)
	}
	defer scriptfile.Close()

	hash := hasht.New()
	io.Copy(hash, scriptfile) //feed data to hash func
	hashstr := fmt.Sprintf("%x", hash.Sum()) //get hash as hex string

	table, err := readTable(dbfilename) //get our data ready
	if err != nil {
		log.Fatalln(err)
	}

	metadata, ok := table[hashstr] //look for record of scriptfile
	if !ok {
		metadata, err = compile(scriptfile)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		if _, err = os.Stat(path.Join(storedir, hashstr)); err != nil {
			metadata, err = compile(scriptfile)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			metadata.lastused, _, err = os.Time()
			if err != nil {
				log.Println(err)
				metadata.lastused = ^(1 << 63) //set time to latest possible
			}
			metadata.filename = scriptname
		}
	}
	table[hashstr] = metadata
	if err := writeTable(table, dbfilename); err != nil {
		log.Println(err)
	}

	os.Exec(path.Join(storedir, hashstr), os.Args[1:], os.Environ())
}
