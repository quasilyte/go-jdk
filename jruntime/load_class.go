package jruntime

import (
	"github.com/quasilyte/GopherJRE/ir"
)

func loadClass(vm *VM, irclass *ir.Class) (*Class, error) {
	methods := make([]Method, len(irclass.Methods))
	for i := range methods {
		m := &methods[i]
		m.Name = irclass.Methods[i].Name
	}

	class := &Class{
		Name:    irclass.Name,
		Methods: methods,
	}
	return class, nil
}
