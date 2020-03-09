package jruntime

// jcall runs method code inside env context.
//
// Code is a pointer to a beginning of a method machine code.
func jcall(e *Env, code *byte)
