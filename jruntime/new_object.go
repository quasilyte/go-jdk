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
	var data *int32
	if length != 0 {
		elems := make([]int32, length)
		data = &elems[0]
	}
	return (*Object)(unsafe.Pointer(&IntArrayObject{
		Info: &IntArrayInfo,
		Data: data,
		Len:  length,
	}))
}
