package jruntime

import (
	"unsafe"

	"github.com/quasilyte/go-jdk/goreflect"
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
	Info *ObjectInfo
	Data *int32
	Len  int32
}

func (o *IntArrayObject) AsSlice() []int32 {
	header := goreflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(o.Data)),
		Len:  int(o.Len),
		Cap:  int(o.Len),
	}
	return *(*[]int32)(unsafe.Pointer(&header))
}

func (o *Object) AsIntArray() *IntArrayObject {
	return (*IntArrayObject)(unsafe.Pointer(o))
}
