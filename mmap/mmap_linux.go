package mmap

import (
	"syscall"
	"unsafe"

	"github.com/quasilyte/go-jdk/goreflect"
)

func mmapExecutable(length int) (Descriptor, []byte, error) {
	prot := syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
	flags := syscall.MAP_PRIVATE | syscall.MAP_ANON
	return mmapLinux(0, length, prot, flags, 0, 0)
}

func munmap(d Descriptor) error {
	return munmapLinux(d)
}

func mmapLinux(addr uintptr, length, prot, flags int, fd uintptr, offset int64) (Descriptor, []byte, error) {
	// void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset);
	ptr, _, err := syscall.Syscall6(
		syscall.SYS_MMAP,
		addr,
		uintptr(length),
		uintptr(prot),
		uintptr(flags),
		fd,
		uintptr(offset))
	d := Descriptor{Addr: ptr, Size: length}
	if err != 0 {
		return d, nil, err
	}
	slice := *(*[]byte)(unsafe.Pointer(&goreflect.SliceHeader{
		Data: ptr,
		Len:  length,
		Cap:  length,
	}))
	return d, slice, nil
}

func munmapLinux(d Descriptor) error {
	// int munmap(void *addr, size_t length);
	_, _, err := syscall.Syscall(
		syscall.SYS_MUNMAP,
		d.Addr,
		uintptr(d.Size),
		0)
	if err != 0 {
		return err
	}
	return nil
}
