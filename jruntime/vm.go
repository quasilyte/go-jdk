package jruntime

import (
	"github.com/quasilyte/GopherJRE/ir"
)

type VM struct {
	// This is a stub for the future.
	//
	// VM is supposed to store all JIT'ed code,
	// loaded classes and their state.
	// Every execution unit (Env) has a link to this object.
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) LoadClass(c *ir.Class) (*Class, error) {
	return loadClass(vm, c)
}
