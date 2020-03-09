package jruntime

import (
	"unsafe"

	"github.com/quasilyte/go-jdk/jit"
)

func BindFuncs(ctx *jit.Context) {
	ctx.Funcs.Jcall = funcAddr(jcall)
	ctx.Funcs.NewIntArray = funcAddr(NewIntArray)
}

// funcAddr returns function value fn executable code address.
func funcAddr(fn interface{}) uint32 {
	// emptyInterface is the header for an interface{} value.
	type emptyInterface struct {
		typ   uintptr
		value *uintptr
	}
	e := (*emptyInterface)(unsafe.Pointer(&fn))
	addr := *e.value
	if addr > 0xffffffff {
		panic("func addr does not fit in 32-bit")
	}
	return uint32(addr)
}
