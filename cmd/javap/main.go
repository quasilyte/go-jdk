package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/quasilyte/GopherJRE/cmd/internal/cmdutil"
	"github.com/quasilyte/GopherJRE/irgen"
	"github.com/quasilyte/GopherJRE/javap"
	"github.com/quasilyte/GopherJRE/jclass"
	"github.com/quasilyte/GopherJRE/jruntime"
)

func main() {
	var args arguments
	flag.StringVar(&args.format, "format", "raw",
		`output format: raw, ir or asm`)
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
	jf, err := cmdutil.DecodeClassFile(filename)
	if err != nil {
		return fmt.Errorf("decode error: %v", err)
	}

	switch args.format {
	case "raw":
		javap.Fprint(os.Stdout, jf)
		return nil
	case "ir":
		return printFileIR(jf)
	case "asm":
		return printFileAsm(jf)
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
		fmt.Printf("  method %s:\n", m.Name)
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

func printFileAsm(jf *jclass.File) error {
	irclass, err := irgen.ConvertClass(jf)
	if err != nil {
		return err
	}

	compiler := jruntime.NewCompiler(runtime.GOARCH)
	if compiler == nil {
		return fmt.Errorf("%s arch is not supported", runtime.GOARCH)
	}
	vm := jruntime.NewVM(compiler)
	defer vm.Close()
	c, err := vm.LoadClass(irclass)
	if err != nil {
		return err
	}

	for _, m := range c.Methods {
		fmt.Printf("  method %s:\n", m.Name)
		fmt.Printf("    %x\n", m.Code)
	}

	return nil
}
