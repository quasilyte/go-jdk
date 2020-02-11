package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/quasilyte/GopherJRE/javap"
	"github.com/quasilyte/GopherJRE/jclass"
)

func main() {
	flag.Parse()

	filenames := flag.Args()
	for _, filename := range filenames {
		if err := printFile(filename); err != nil {
			log.Printf("%s: %v", filename, err)
		}
	}
}

func printFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var dec jclass.Decoder
	jf, err := dec.Decode(f)
	if err != nil {
		return fmt.Errorf("decode error: %v", err)
	}

	javap.Fprint(os.Stdout, jf)
	return nil
}
