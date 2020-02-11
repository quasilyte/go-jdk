package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/quasilyte/GopherJRE/irgen"
	"github.com/quasilyte/GopherJRE/javap"
	"github.com/quasilyte/GopherJRE/jclass"
)

func main() {
	var args arguments
	flag.StringVar(&args.format, "format", "raw",
		`output format: raw or ir`)
	flag.Parse()

	filenames := flag.Args()
	for _, filename := range filenames {
		if err := printFile(&args, filename); err != nil {
			log.Printf("%s: %v", filename, err)
		}
	}
}

type arguments struct {
	format string
}

func printFile(args *arguments, filename string) error {
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

	switch args.format {
	case "raw":
		javap.Fprint(os.Stdout, jf)
		return nil
	case "ir":
		return printFileIR(jf)
	default:
		return fmt.Errorf("unknown format: %s", args.format)
	}
}

func printFileIR(jf *jclass.File) error {
	c, err := irgen.ConvertClass(jf)
	if err != nil {
		return err
	}

	fmt.Printf("name=%q pkg=%q\n", c.Name, c.PkgName)
	for _, m := range c.Methods {
		fmt.Printf("  method=%q\n", m.Name)
		blockIndex := -1
		for i, inst := range m.Code {
			if inst.Flags.IsBlockLead() {
				blockIndex++
			}
			fmt.Printf("        b%d %3d: %s\n", blockIndex, i, inst)
		}
	}

	return nil
}
