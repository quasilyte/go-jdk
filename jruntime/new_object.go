package jruntime

import (
	"unsafe"
)

var IntArrayInfo = ObjectInfo{
	Kind: KindArray,
	Size: 4,
}

func NewIntArray(env *Env, length int32) *Object {
	env.trackAllocation(int64(length)*4 + 4)
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

func NewObject4(env *Env) *Object {
	env.trackAllocation(4)
	return (*Object)(unsafe.Pointer(&object4{}))
}

func NewObject8(env *Env) *Object {
	env.trackAllocation(8)
	return (*Object)(unsafe.Pointer(&object8{}))
}
