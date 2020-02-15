package mmap

type Descriptor struct {
	Addr uintptr
	Size int
}

// Executable allocates a mapped memory buffer that can be
// used to store runnable machine code.
//
// Descriptor should be used to free (unmap) the memory.
func Executable(length int) (Descriptor, []byte, error) {
	return executable(length)
}

// Free unmaps the memory that is associated with the given descriptor.
func Free(d Descriptor) error {
	return free(d)
}
