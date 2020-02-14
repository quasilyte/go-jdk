package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/quasilyte/GopherJRE/cmd/internal/cmdutil"
	"github.com/quasilyte/GopherJRE/irgen"
	"github.com/quasilyte/GopherJRE/jruntime"
)

func main() {
	log.SetFlags(0)

	var args arguments
	flag.StringVar(&args.classFile, "class", "",
		`path to a class file`)
	flag.StringVar(&args.methodName, "method", "main",
		`class static method name`)
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

	vm := jruntime.NewVM()

	classfile, err := cmdutil.DecodeClassFile(args.classFile)
	if err != nil {
		log.Fatalf("decode class file: %v", err)
	}
	irclass, err := irgen.ConvertClass(classfile)
	if err != nil {
		log.Fatalf("generate ir: %v", err)
	}
	compileStart := time.Now()
	class, err := vm.LoadClass(irclass)
	compileTime := time.Since(compileStart)
	if err != nil {
		log.Fatalf("load class: %v", err)
	}

	method := class.FindMethod(args.methodName)
	if method == nil {
		log.Fatalf("method %s.%s not found", class.Name, args.methodName)
	}

	env := jruntime.NewEnv(vm, &jruntime.EnvConfig{})
	for i, arg := range args.methodArgs {
		v, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			log.Fatalf("method arg[%d]: only int args are supported, found %s", i, arg)
		}
		env.IntArg(i, v)
	}

	callStart := time.Now()
	result, err := env.IntCall(method)
	callTime := time.Since(callStart)
	if err != nil {
		log.Fatalf("call returned error: %v", err)
	}

	argsString := strings.Join(args.methodArgs, ", ")
	log.Printf("%s.%s(%s) => %d\n", class.Name, args.methodName, argsString, result)
	if args.verbose {
		log.Println("-- verbose output --")
		log.Printf("method machine code: %x\n", method.Code)
		log.Printf("class load time:  %.8fs (%d ns)\n",
			compileTime.Seconds(), compileTime.Nanoseconds())
		log.Printf("method call time: %.8fs (%d ns)\n",
			callTime.Seconds(), callTime.Nanoseconds())
	}
}

type arguments struct {
	classFile  string
	methodName string
	methodArgs []string
	verbose    bool
}
