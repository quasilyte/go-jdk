package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/quasilyte/go-jdk/irgen"
	"github.com/quasilyte/go-jdk/jit"
	"github.com/quasilyte/go-jdk/jruntime"
	"github.com/quasilyte/go-jdk/loader"
	"github.com/quasilyte/go-jdk/vmdat"
)

func runMain() error {
	var cmd runCommand
	flag.StringVar(&cmd.classFile, "class", "",
		`path to a class file`)
	flag.StringVar(&cmd.methodName, "method", "main",
		`class static method name`)
	flag.StringVar(&cmd.methodSignature, "signature", "",
		`method signature which is required if method is overloaded`)
	flag.BoolVar(&cmd.verbose, "v", false,
		`verbose output mode`)
	flag.StringVar(&cmd.classPath, "cp", "",
		`class path to use`)
	flag.Parse()
	cmd.methodArgs = flag.Args()

	if cmd.classFile == "" {
		return errors.New("-class argument can't be empty")
	}
	if cmd.methodName == "" {
		return errors.New("-method argument can't be empty")
	}

	return cmd.run()
}

type runCommand struct {
	classFile       string
	methodName      string
	methodSignature string
	methodArgs      []string
	classPath       string
	verbose         bool
}

func (cmd *runCommand) run() error {
	vm, err := jruntime.OpenVM(runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("open VM: %v", err)
	}
	defer vm.Close()

	compileStart := time.Now()
	class, err := cmd.loadAndCompileClass(vm, cmd.classFile)
	compileTime := time.Since(compileStart)
	if err != nil {
		return fmt.Errorf("load class: %v", err)
	}

	method := class.FindMethod(cmd.methodName, cmd.methodSignature)
	if method == nil {
		return fmt.Errorf("method %s.%s not found", class.Name, cmd.methodName)
	}

	env := jruntime.NewEnv(vm, &jruntime.EnvConfig{})
	for i, arg := range cmd.methodArgs {
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

	argsString := strings.Join(cmd.methodArgs, ", ")
	log.Printf("%s.%s(%s) => %d\n", class.Name, method.Name, argsString, result)
	if cmd.verbose {
		log.Println("-- verbose output --")
		log.Printf("method machine code: %x\n", method.Code)
		log.Printf("class load time:  %.8fs (%d ns)\n",
			compileTime.Seconds(), compileTime.Nanoseconds())
		log.Printf("method call time: %.8fs (%d ns)\n",
			callTime.Seconds(), callTime.Nanoseconds())
	}

	return nil
}

func (cmd *runCommand) loadAndCompileClass(vm *jruntime.VM, filename string) (*vmdat.Class, error) {
	toCompile, err := loader.LoadClass(&vm.State, filename, &loader.Config{
		ClassPath: []string{cmd.classPath},
	})
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
