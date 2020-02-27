package jit

import (
	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/mmap"
	"github.com/quasilyte/go-jdk/vmdat"
)

type Context struct {
	State *vmdat.State
	Mmap  *mmap.Manager

	JcallAddr uint32
}

// Compiler is used by a VM to generate machine code for class methods.
type Compiler interface {
	Compile(Context, []*ir.Package) error
}
