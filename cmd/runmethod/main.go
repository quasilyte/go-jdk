package main

import (
	"flag"
	"log"
	"os"

	"github.com/quasilyte/GopherJRE/irgen"
	"github.com/quasilyte/GopherJRE/jclass"
	"github.com/quasilyte/GopherJRE/jruntime"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("runmethod: ")

	var args arguments
	flag.StringVar(&args.classFile, "class", "",
		`path to a class file`)
	flag.StringVar(&args.methodName, "method", "main",
		`class static method name`)
	flag.Parse()
	args.methodArgs = flag.Args()

	if args.classFile == "" {
		log.Fatalf("-class argument can't be empty")
	}
	if args.methodName == "" {
		log.Fatalf("-method argument can't be empty")
	}

	vm := jruntime.NewVM()

	classfile, err := decodeClassFile(args.classFile)
	if err != nil {
		log.Fatalf("decode class file: %v", err)
	}
	irclass, err := irgen.ConvertClass(classfile)
	if err != nil {
		log.Fatalf("generate ir: %v", err)
	}
	class, err := vm.LoadClass(irclass)
	if err != nil {
		log.Fatalf("load class: %v", err)
	}

	method := class.FindMethod(args.methodName)
	if method == nil {
		log.Fatalf("method %s.%s not found", class.Name, args.methodName)
	}

	env := jruntime.NewEnv(vm, &jruntime.EnvConfig{})
	if err := env.Call(method); err != nil {
		log.Fatalf("call returned error: %v", err)
	}

	log.Printf("success")
}

type arguments struct {
	classFile  string
	methodName string
	methodArgs []string
}

func decodeClassFile(filename string) (*jclass.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var dec jclass.Decoder
	return dec.Decode(f)
}
