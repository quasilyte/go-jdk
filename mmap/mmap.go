package mmap

// Executable allocates a mapped memory buffer that can be
// used to store runnable machine code.
func Executable(length int) ([]byte, error) {
	return executable(length)
}
