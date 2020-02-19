package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/quasilyte/GopherJRE/irgen"
	"github.com/quasilyte/GopherJRE/jit"
	"github.com/quasilyte/GopherJRE/jruntime"
	"github.com/quasilyte/GopherJRE/loader"
	"github.com/quasilyte/GopherJRE/vmdat"
)

func main() {
	log.SetFlags(0)

	var args arguments
	flag.StringVar(&args.classFile, "class", "",
		`path to a class file`)
	flag.StringVar(&args.methodName, "method", "main",
		`class static method name`)
	flag.StringVar(&args.methodSignature, "signature", "",
		`method signature which is required if method is overloaded`)
	flag.BoolVar(&args.verbose, "v", false,
		`verbose output mode`)
	flag.Parse()
	args.methodArgs = flag.Args()

	if args.classFile == "" {
		log.Fatalf("-class argument can't be empty")
	}
	if args.methodName == "" {
		log.Fatalf("-method argument can't be empty")
	}

	if err := mainNoExit(&args); err != nil {
		log.Fatal(err)
	}
}

type arguments struct {
	classFile       string
	methodName      string
	methodSignature string
	methodArgs      []string
	verbose         bool
}

func mainNoExit(args *arguments) error {
	vm, err := jruntime.OpenVM(runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("open VM: %v", err)
	}
	defer vm.Close()

	compileStart := time.Now()
	class, err := loadAndCompileClass(vm, args.classFile)
	compileTime := time.Since(compileStart)
	if err != nil {
		return fmt.Errorf("load class: %v", err)
	}

	method := class.FindMethod(args.methodName, args.methodSignature)
	if method == nil {
		return fmt.Errorf("method %s.%s not found", class.Name, args.methodName)
	}

	env := jruntime.NewEnv(vm, &jruntime.EnvConfig{})
	for i, arg := range args.methodArgs {
		v, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("method arg[%d]: only int args are supported, found %s", i, arg)
		}
		env.IntArg(i, v)
	}

	callStart := time.Now()
	result, err := env.IntCall(method)
	callTime := time.Since(callStart)
	if err != nil {
		return fmt.Errorf("call returned error: %v", err)
	}

	argsString := strings.Join(args.methodArgs, ", ")
	log.Printf("%s.%s(%s) => %d\n", class.Name, method.Name, argsString, result)
	if args.verbose {
		log.Println("-- verbose output --")
		log.Printf("method machine code: %x\n", method.Code)
		log.Printf("class load time:  %.8fs (%d ns)\n",
			compileTime.Seconds(), compileTime.Nanoseconds())
		log.Printf("method call time: %.8fs (%d ns)\n",
			callTime.Seconds(), callTime.Nanoseconds())
	}

	return nil
}

func loadAndCompileClass(vm *jruntime.VM, filename string) (*vmdat.Class, error) {
	toCompile, err := loader.LoadClass(&vm.State, filename, &loader.Config{})
	if err != nil {
		return nil, fmt.Errorf("load class: %v", err)
	}
	if err := irgen.Generate(&vm.State, toCompile); err != nil {
		return nil, fmt.Errorf("irgen: %v", err)
	}
	ctx := jit.Context{
		Mmap:  &vm.Mmap,
		State: &vm.State,
	}
	if err := vm.Compiler.Compile(ctx, toCompile); err != nil {
		return nil, fmt.Errorf("compile: %v", err)
	}
	className := strings.TrimSuffix(filepath.Base(filename), ".class")
	return toCompile[0].Out.FindClass(className), nil
}
