package jruntime

import (
	"fmt"

	"github.com/quasilyte/GopherJRE/codegen"
	"github.com/quasilyte/GopherJRE/ir"
	"github.com/quasilyte/GopherJRE/mmap"
)

type VM struct {
	compiler codegen.Compiler

	mapped []mmap.Descriptor
}

func NewVM(compiler codegen.Compiler) *VM {
	return &VM{compiler: compiler}
}

// Close instructs the VM to free all associated resources.
func (vm *VM) Close() error {
	for _, d := range vm.mapped {
		if err := mmap.Free(d); err != nil {
			return fmt.Errorf("free mapped memory: %v", err)
		}
	}
	return nil
}

// LoadClass attempts to load a given class.
// All dependencies required by c should be already loaded.
func (vm *VM) LoadClass(c *ir.Class) (*Class, error) {
	methods := make([]Method, len(c.Methods))

	compiled, err := vm.compiler.CompileMethods(c.Methods)
	if err != nil {
		return nil, err
	}
	for i, code := range compiled.Methods {
		m := &methods[i]
		m.Name = c.Methods[i].Name

		if len(code) == 0 {
			continue
		}
		executable, err := vm.getCodeBuf(len(code))
		if err != nil {
			return nil, fmt.Errorf("mmap(%d): %v", len(code), err)
		}
		copy(executable, code)
		m.Code = executable
	}

	class := &Class{Name: c.Name, Methods: methods}
	return class, nil
}

func (vm *VM) getCodeBuf(length int) ([]byte, error) {
	// TODO(quasilyte): re-use mapped regions.
	d, buf, err := mmap.Executable(length)
	if err != nil {
		return nil, err
	}
	vm.mapped = append(vm.mapped, d)
	return buf, nil
}
