package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/quasilyte/go-jdk/cmd/internal/cmdutil"
	"github.com/quasilyte/go-jdk/jdeps"
)

func main() {
	var args arguments
	flag.StringVar(&args.filename, "filename", "", `path to class file`)
	flag.Parse()

	file, err := cmdutil.DecodeClassFile(args.filename)

	if err != nil {
		log.Fatal(err)
	}

	deps := jdeps.ClassDependencies(file)

	fmt.Printf("%v", deps)
}

type arguments struct {
	filename string
}
