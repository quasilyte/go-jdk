package jruntime

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/quasilyte/GopherJRE/mmap"
)

// TODO(quasilyte): can we make this test less arch-specific?

func TestJcall(t *testing.T) {
	var env Env
	stack := make([]int64, 16)
	env.stack = (*byte)(unsafe.Pointer(&stack[0]))

	funcCode := []byte{
		// MOVQ $777, (SI)
		0x48, 0xc7, 0x06, 0x09, 0x03, 0x00, 0x00,
		// MOVL $1, AX
		0xb8, 0x01, 0x00, 0x00, 0x00,
		// JMP -8(SI)
		0xff, 0x66, 0xf8,
	}
	d, executable, err := mmap.Executable(len(funcCode))
	if err != nil {
		t.Errorf("mmap executable: %v", err)
	}
	copy(executable, funcCode)
	defer munmap(t, d)

	checkStack := func() {
		if stack[0] != 1 {
			t.Errorf("stack[0] mismatch:\nhave: %d\nwant: 1", stack[0])
		}
		if stack[1] != 777 {
			t.Errorf("stack[1] mismatch:\nhave: %d\nwant: 777", stack[1])
		}
		stack[0] = 0
		stack[1] = 0
	}

	// Call in a same goroutine, on a same frame.
	jcall(&env, &executable[0])
	checkStack()

	// Should not panic nor corrupt the stack.
	for i := 0; i < 10; i++ {
		nestedCall(&env, executable)
		checkStack()
	}

	// Now test it inside a new goroutine.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		jcall(&env, &executable[0])
		checkStack()
		wg.Done()
	}()
	wg.Wait()
}

//go:noinline
func nestedCall(env *Env, code []byte) {
	jcall(env, &code[0])
}

func BenchmarkJcall(b *testing.B) {
	var env Env
	stack := make([]int64, 16)
	env.stack = (*byte)(unsafe.Pointer(&stack[0]))

	funcCode := []byte{
		// JMP -8(SI)
		0xff, 0x66, 0xf8,
	}
	d, executable, err := mmap.Executable(len(funcCode))
	if err != nil {
		b.Errorf("mmap executable: %v", err)
	}
	copy(executable, funcCode)
	defer munmap(b, d)
	executablePtr := &executable[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jcall(&env, executablePtr)
	}
	b.StopTimer()
}

func munmap(tb testing.TB, d mmap.Descriptor) {
	if err := mmap.Free(d); err != nil {
		tb.Errorf("mmap free: %v", err)
	}
}
