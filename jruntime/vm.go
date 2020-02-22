package jruntime

import (
	"fmt"

	"github.com/quasilyte/go-jdk/jit"
	"github.com/quasilyte/go-jdk/jit/compiler/x64"
	"github.com/quasilyte/go-jdk/mmap"
	"github.com/quasilyte/go-jdk/vmdat"
)

type VM struct {
	State    vmdat.State
	Mmap     mmap.Manager
	Compiler jit.Compiler
}

func OpenVM(arch string) (*VM, error) {
	var vm VM
	switch arch {
	case "amd64":
		vm.Compiler = x64.NewCompiler()
	default:
		return nil, fmt.Errorf("arch %s is not supported", arch)
	}
	vm.State.Init()
	return &vm, nil
}

// Close instructs the VM to free all associated resources.
func (vm *VM) Close() error {
	if err := vm.Mmap.Close(); err != nil {
		return err
	}
	return nil
}
