package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/quasilyte/GopherJRE/cmd/internal/cmdutil"
	"github.com/quasilyte/GopherJRE/ir"
	"github.com/quasilyte/GopherJRE/irgen"
	"github.com/quasilyte/GopherJRE/javap"
	"github.com/quasilyte/GopherJRE/jit"
	"github.com/quasilyte/GopherJRE/jruntime"
	"github.com/quasilyte/GopherJRE/loader"
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
	if args.format == "raw" {
		jf, err := cmdutil.DecodeClassFile(filename)
		if err != nil {
			return fmt.Errorf("decode error: %v", err)
		}
		javap.Fprint(os.Stdout, jf)
		return nil
	}

	if args.format != "ir" {
		return fmt.Errorf("unknown format: %s", args.format)
	}

	vm, err := jruntime.OpenVM(runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("open VM: %v", err)
	}
	defer vm.Close()

	toCompile, err := loader.LoadClass(&vm.State, filename, &loader.Config{})
	if err != nil {
		return fmt.Errorf("load class: %v", err)
	}
	if err := irgen.Generate(&vm.State, toCompile); err != nil {
		return fmt.Errorf("irgen: %v", err)
	}
	ctx := jit.Context{
		Mmap:  &vm.Mmap,
		State: &vm.State,
	}
	if err := vm.Compiler.Compile(ctx, toCompile); err != nil {
		return fmt.Errorf("compile: %v", err)
	}
	pkg := toCompile[0]
	className := strings.TrimSuffix(filepath.Base(filename), ".class")
	var class *ir.Class
	for i := range pkg.Classes {
		if pkg.Classes[i].Out.Name == className {
			class = &pkg.Classes[i]
			break
		}
	}

	fmt.Printf("class %q\n", class.Name)
	for i := range class.Methods {
		m := &class.Methods[i]
		fmt.Printf("  method %s (slots: %d):\n", m.Out.Name, m.Out.FrameSlots)
		blockIndex := -1
		for i, inst := range m.Code {
			if inst.Flags.IsBlockLead() {
				blockIndex++
			}
			fmt.Printf("        b%d %3d: %s\n", blockIndex, i, inst)
		}
		code := m.Out.Code
		for len(code) != 0 {
			length := 32
			if len(code) < length {
				length = len(code)
			}
			fmt.Printf("        %x\n", code[:length])
			code = code[length:]
		}
	}

	return nil
}
