package jit

import (
	"github.com/quasilyte/GopherJRE/ir"
	"github.com/quasilyte/GopherJRE/mmap"
	"github.com/quasilyte/GopherJRE/vmdat"
)

type Context struct {
	State *vmdat.State
	Mmap  *mmap.Manager
}

// Compiler is used by a VM to generate machine code for class methods.
type Compiler interface {
	Compile(Context, []*ir.Package) error
}
