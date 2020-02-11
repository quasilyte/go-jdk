package jruntime

import (
	"sync/atomic"
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

	// StackMemory is a memory that can be used by a running program.
	//
	// Stack is used to allocate local stack slots, call frames, etc.
	// Setting lower value effectively limits how deep the recursion can go.
	// If there is not enough stack memory for another function call,
	// VM execution is terminated.
	//
	// Slice ownership is passed into Env, it should not be
	// modified outside of it afterwards.
	//
	// nil will result in a new memory allocation of default size.
	StackMemory []byte
}

// NewEnv returns configured execution unit for the VM.
func NewEnv(vm *VM, cfg *EnvConfig) *Env {
	const megabyte = 1024 * 1024

	var env Env
	env.vm = vm

	if cfg.StackMemory == nil {
		stack := make([]byte, megabyte/4)
		env.stack = &stack[0]
	} else {
		env.stack = &cfg.StackMemory[0]
	}

	if cfg.AllocBytesLimit == 0 {
		env.allocBytesLimit = megabyte
	}

	return &env
}

// Env is context-like structure for VM code execution.
type Env struct {
	envFixed // Should be the first struct member

	allocBytesLimit int64
}

type envFixed struct {
	// Note: please don't re-arrange members of this struct
	// as they are sometimes accessed via computed offsets manually.
	// Everything that doesn't need to be aligned carefully can go
	// into Env struct itself instead.

	allocBytesLeft int64

	stack *byte

	vm *VM
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
