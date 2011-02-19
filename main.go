package main

import (
	"os"
	"io"
	"fmt"
	"log"
	hasht "crypto/sha256" //hash type, changeable
)

func main() {

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
		table[hashstr] = metadata
	} else {
		if !os.??fileexists(storedir + hashstr) {
			metadata, err = compile(scriptfile)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			metadata.lastused = os.Time()
			metadata.filename = scriptname
		}
		table[hashstr] = metadata
	}
	if err := writeTable(table, dbfilename); err != nil {
		log.Println(err)
	}

	os.Exec(hashstr, os.Args[1:], os.Environ())
}
