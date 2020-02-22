package main

import (
	"log"
)

func main() {
	log.SetFlags(0)
	dispatchCommand(subCommands)
}

var subCommands = []*subCommand{
	{
		Main:  jdepsMain,
		Name:  "jdeps",
		Short: "print JVM class file dependencies",
		Examples: []string{
			"go-jdk jdeps -help",
			"go-jdk jdeps /path/to/File.class",
		},
	},

	{
		Main:  javapMain,
		Name:  "javap",
		Short: "disassemble JVM class or print its IR representation",
		Examples: []string{
			"go-jdk javap -help",
			"go-jdk javap /path/to/File.class",
			"go-jdk javap -format=ir File.class",
		},
	},

	{
		Main:  runMain,
		Name:  "run",
		Short: "load and run JVM class",
		Examples: []string{
			"go-jdk run -help",
			"go-jdk run /path/to/File.class",
			"go-jdk run -method add Arith.class 1 2",
		},
	},

	{
		Main:     versionMain,
		Name:     "version",
		Short:    "print go-jdk version",
		Examples: []string{"go-jdk version"},
	},
}
