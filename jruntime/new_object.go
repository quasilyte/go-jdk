package jruntime

import (
	"unsafe"
)

var IntArrayInfo = ObjectInfo{
	Kind: KindArray,
	Size: 4,
}

func NewIntArray(env *Env, length int32) *Object {
	env.trackAllocation(int64(length)*4 + int64(unsafe.Sizeof(IntArrayObject{})))
	elems := make([]int32, length)
	return (*Object)(unsafe.Pointer(&IntArrayObject{
		Info:  &IntArrayInfo,
		Elems: &elems[0],
		Len:   length,
	}))
}
