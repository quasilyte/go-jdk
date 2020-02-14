package mmap

import (
	"reflect"
	"syscall"
	"unsafe"
)

func executable(length int) ([]byte, error) {
	prot := syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
	flags := syscall.MAP_PRIVATE | syscall.MAP_ANON
	return mmap(0, length, prot, flags, 0, 0)
}

func mmap(addr uintptr, length, prot, flags int, fd uintptr, offset int64) ([]byte, error) {
	ptr, _, err := syscall.Syscall6(
		syscall.SYS_MMAP,
		addr,
		uintptr(length),
		uintptr(prot),
		uintptr(flags),
		fd,
		uintptr(offset))
	if err != 0 {
		return nil, err
	}
	slice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: ptr,
		Len:  length,
		Cap:  length,
	}))
	return slice, nil
}
