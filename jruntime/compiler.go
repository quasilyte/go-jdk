package jruntime

import (
	"github.com/quasilyte/GopherJRE/codegen"
	"github.com/quasilyte/GopherJRE/compiler/x64"
)

// NewCompiler returns a compiler for a host machine architecture.
// If there is no compiler for that architecture, nil is returned.
func NewCompiler(arch string) codegen.Compiler {
	switch arch {
	case "amd64":
		return x64.NewCompiler()
	default:
		return nil
	}
}
