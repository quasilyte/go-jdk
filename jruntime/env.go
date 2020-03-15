package jruntime

import (
	"sync/atomic"
	"unsafe"

	"github.com/quasilyte/go-jdk/vmdat"
)

// TODO(quasilyte): decide whether Env is thread-safe and in what contexts.
// This can wait for until we run VM inside a concurrent program.

// TODO(quasilyte): define "VM execution is terminated" more precisely.

// EnvConfig describes VM execution unit settings.
type EnvConfig struct {
	// AllocBytesLimit describes max alloc bytes limit for a single invocation.
	//
	// When VM function is called, alloc counter is set to this value.
	// Every time a new object is allocated, it's size is subtracted from that counter.
	// If it becomes negative, VM execution is terminated.
	//
	// Zero value means "default value" which is big enough for most use cases.
	AllocBytesLimit int64

	// StackMemory is a stack size in bytes that can be used by a running program.
	//
	// Stack is used to allocate local stack slots, call frames, etc.
	// Setting lower value effectively limits how deep the recursion can go.
	// If there is not enough stack memory for another function call,
	// VM execution is terminated.
	//
	// Value of 0 will result in in a default stack size allocation.
	StackMemory int64
}

// NewEnv returns configured execution unit for the VM.
func NewEnv(vm *VM, cfg *EnvConfig) *Env {
	const megabyte = 1024 * 1024

	var env Env
	env.vm = vm

	stackMemory := cfg.StackMemory
	if stackMemory == 0 {
		stackMemory = megabyte
	}
	env.slots = make([]stackSlot, stackMemory)
	env.stack = &env.slots[0]

	env.allocBytesLimit = cfg.AllocBytesLimit
	if env.allocBytesLimit == 0 {
		env.allocBytesLimit = megabyte
	}

	return &env
}

// Env is context-like structure for VM code execution.
type Env struct {
	envFixed // Should be the first struct member

	allocBytesLimit int64

	slots []stackSlot
}

type stackSlot struct {
	scalar int64
	ptr    *Object
}

type envFixed struct {
	// Note: please don't re-arrange members of this struct
	// as they are sometimes accessed via computed offsets manually.
	// Everything that doesn't need to be aligned carefully can go
	// into Env struct itself instead.
	// FIXME: this layout implies 64-bit platform.

	allocBytesLeft int64      // offset=0
	stack          *stackSlot // offset=8
	_              uint64     // offset=16 (reserved)
	tmp            uint64     // offset=24

	vm *VM
}

func (env *Env) GC() {
	for _, slot := range env.slots {
		slot.ptr = nil
	}
}

func (env *Env) IntCall(m *vmdat.Method) (int64, error) {
	env.allocBytesLeft = env.allocBytesLimit
	jcallScalar(env, &m.Code[0])
	return env.stack.scalar, nil
}

func (env *Env) IntArg(i int, v int64) {
	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(env.stack)) + uintptr(i*16) + 16)
	(*stackSlot)(ptr).scalar = v
}

// trackAllocations checks whether we can allocate size bytes.
// If memory limit is reached, it panics.
func (env *Env) trackAllocation(size int64) {
	n := atomic.AddInt64(&env.allocBytesLeft, -size)
	if n < 0 {
		// TODO(quasilyte): provide JVM stack trace info.
		panic("allocations limit reached")
	}
}
