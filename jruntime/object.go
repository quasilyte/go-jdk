package jruntime

import (
	"unsafe"
)

type Object struct {
	Info    *ObjectInfo
	Ptrdata unsafe.Pointer
}

type ObjectInfo struct {
	Kind ObjectKind
	Size int
}

type ObjectKind int

const (
	KindArray ObjectKind = iota
	KindObject
)

type IntArrayObject struct {
	Info  *ObjectInfo
	Elems *int32
	Len   int32
}

func (o *Object) AsIntArray() *IntArrayObject {
	return (*IntArrayObject)(unsafe.Pointer(o))
}
