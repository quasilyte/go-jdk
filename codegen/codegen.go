package codegen

import (
	"github.com/quasilyte/GopherJRE/ir"
)

// Compiler is used by a VM to load classes.
//
// If compiler also implements io.Closer, VM will call that
// method during its own Close().
type Compiler interface {
	// CompileMethods returns native machine code for class methods.
	//
	// Only the first compilation error is returned.
	CompileMethods([]ir.Method) (Result, error)
}

// Result represents methods compilation result.
type Result struct {
	Methods [][]byte
	Relocs  []Relocation
}

// Relocation is a handle to a symbol that needs to be resolved.
type Relocation struct {
	Fixup *uint64
}
