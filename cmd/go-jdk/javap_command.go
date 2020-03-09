package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/quasilyte/go-jdk/cmd/internal/cmdutil"
	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/irfmt"
	"github.com/quasilyte/go-jdk/irgen"
	"github.com/quasilyte/go-jdk/javap"
	"github.com/quasilyte/go-jdk/jit"
	"github.com/quasilyte/go-jdk/jruntime"
	"github.com/quasilyte/go-jdk/loader"
)

func javapMain() error {
	var cmd javapCommand
	flag.StringVar(&cmd.format, "format", "raw",
		`output format: raw or ir`)
	flag.StringVar(&cmd.classPath, "cp", "",
		`class path to use`)
	flag.Parse()

	filenames := flag.Args()
	for _, filename := range filenames {
		if err := cmd.printFile(filename); err != nil {
			return fmt.Errorf("%s: %v", filename, err)
		}
	}

	return nil
}

type javapCommand struct {
	format    string
	classPath string
}

func (cmd *javapCommand) printFile(filename string) error {
	if cmd.format == "raw" {
		jf, err := cmdutil.DecodeClassFile(filename)
		if err != nil {
			return fmt.Errorf("decode error: %v", err)
		}
		javap.Fprint(os.Stdout, jf)
		return nil
	}

	if cmd.format != "ir" {
		return fmt.Errorf("unknown format: %s", cmd.format)
	}

	vm, err := jruntime.OpenVM(runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("open VM: %v", err)
	}
	defer vm.Close()

	toCompile, err := loader.LoadClass(&vm.State, filename, &loader.Config{
		ClassPath: []string{cmd.classPath},
	})
	if err != nil {
		return fmt.Errorf("load class: %v", err)
	}
	if err := irgen.Generate(&vm.State, toCompile); err != nil {
		return fmt.Errorf("irgen: %v", err)
	}
	jitCtx := jit.Context{
		Mmap:  &vm.Mmap,
		State: &vm.State,
	}
	jruntime.BindFuncs(&jitCtx)
	if err := vm.Compiler.Compile(jitCtx, toCompile); err != nil {
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
		fmt.Printf("  method %s (slots=%d):\n", m.Out.Name, m.Out.FrameSlots)
		blockIndex := -1
		for i, inst := range m.Code {
			if inst.Flags.IsBlockLead() {
				blockIndex++
			}
			fmt.Printf("        b%d %3d: %s\n", blockIndex, i, irfmt.Sprint(&vm.State, inst))
		}
		code := m.Out.Code
		var codeChunks []string
		for len(code) != 0 {
			length := 32
			if len(code) < length {
				length = len(code)
			}
			codeChunks = append(codeChunks, fmt.Sprintf("        %x", code[:length]))
			code = code[length:]
		}
		if len(codeChunks) != 0 {
			fmt.Printf("%s (%d bytes)\n", strings.Join(codeChunks, "\n"), len(m.Out.Code))
		}
	}

	return nil
}
