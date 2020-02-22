package main

import (
	"flag"
	"fmt"

	"github.com/quasilyte/go-jdk/cmd/internal/cmdutil"
	"github.com/quasilyte/go-jdk/jdeps"
)

func jdepsMain() error {
	flag.Parse()

	targets := flag.Args()

	for _, target := range targets {
		cf, err := cmdutil.DecodeClassFile(target)
		if err != nil {
			return fmt.Errorf("%s: %v", target, err)
		}
		deps := jdeps.ClassDependencies(cf)
		fmt.Printf("%s depends on:\n", target)
		for _, dep := range deps {
			fmt.Printf("  %s\n", dep)
		}
	}

	return nil
}
