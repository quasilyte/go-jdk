package jruntime

import (
	"unsafe"
)

var IntArrayInfo = ObjectInfo{
	Kind: KindArray,
	Size: 4,
}

func NewIntArray(length int32) *Object {
	elems := make([]int32, length)
	return (*Object)(unsafe.Pointer(&IntArrayObject{
		Info:  &IntArrayInfo,
		Elems: &elems[0],
		Len:   length,
	}))
}

type object4 struct {
	Object
	data [4]byte
}

type object8 struct {
	Object
	data [8]byte
}

func NewObject4() *Object { return (*Object)(unsafe.Pointer(&object4{})) }
func NewObject8() *Object { return (*Object)(unsafe.Pointer(&object8{})) }
