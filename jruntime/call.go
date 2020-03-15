package jruntime

// jcallScalar runs method code inside env context.
//
// Code is a pointer to a beginning of a method machine code.
func jcallScalar(e *Env, code *byte)
