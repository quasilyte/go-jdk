package main

import (
	"log"
	"os"
)

// subCommand is an implementation of a linter sub-command.
type subCommand struct {
	// Main is command entry point.
	Main func() error

	// Name is sub-command name used to execute it.
	Name string

	// Short describes command in one line of text.
	Short string

	// Examples shows one or more sub-command execution examples.
	Examples []string
}

// findSubCommand looks up SubCommand by its name.
// Returns nil if requested command not found.
func findSubCommand(cmdList []*subCommand, name string) *subCommand {
	for _, cmd := range cmdList {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

// printSubCommands prints cmdList info to the logger (usually stderr).
func printSubCommands(cmdList []*subCommand) {
	log.Println("supported sub-commands:")
	for _, cmd := range cmdList {
		log.Printf("\t%s - %s", cmd.Name, cmd.Short)
		for _, ex := range cmd.Examples {
			log.Printf("\t\t$ %s", ex)
		}
	}
}

// dispatchCommand runs sub command out of specified cmdList based on
// the first command line argument.
func dispatchCommand(cmdList []*subCommand) {
	argv := os.Args
	if len(argv) == 2 {
		if argv[1] == "help" || argv[1] == "-help" || argv[1] == "--help" {
			printSubCommands(cmdList)
			os.Exit(0)
		}
	}
	if len(argv) < 2 {
		log.Printf("not enough arguments, expected sub-command name\n\n")
		printSubCommands(cmdList)
		os.Exit(1)
	}

	subIdx := 1 // [0] is program name
	sub := os.Args[subIdx]
	// Erase sub-command argument (index=1) to make it invisible for
	// sub commands themselves.
	os.Args = append(os.Args[:subIdx], os.Args[subIdx+1:]...)

	// Choose and run sub-command main.
	cmd := findSubCommand(cmdList, sub)
	if cmd == nil {
		log.Printf("unknown sub-command: %s\n\n", sub)
		printSubCommands(cmdList)
		os.Exit(1)
	}

	if err := cmd.Main(); err != nil {
		log.Fatal(err)
	}
}
