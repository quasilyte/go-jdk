package goreflect

type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
