package jruntime

import (
	"unsafe"
)

// jcall runs method code inside env context.
//
// Code is a pointer to a beginning of a method machine code.
func jcall(e *Env, code *byte)

// JcallAddr is the address of the jcall function entry point.
var JcallAddr uint32

func init() {
	fn := jcall
	addr := **(**uintptr)(unsafe.Pointer(&fn))
	if addr > 0xffffffff {
		panic("JcallAddr does not fit in 32-bit")
	}
	JcallAddr = uint32(addr)
}
