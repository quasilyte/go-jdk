package jruntime

import (
	"sync/atomic"
)

// TODO(quasilyte): decide whether Env is thread-safe and in what contexts.
// This can wait for until we run VM inside a concurrent program.

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
}

// NewEnv returns configured execution unit for the VM.
func NewEnv(vm *VM, cfg *EnvConfig) *Env {
	const megabyte = 1024 * 1024

	env := Env{vm: vm}
	if cfg.AllocBytesLimit == 0 {
		env.allocBytesLimit = megabyte
	}

	return &env
}

// Env is context-like structure for VM code execution.
type Env struct {
	allocBytesLeft  int64
	allocBytesLimit int64

	vm *VM
}

func (env *Env) trackAllocation(size int64) {
	n := atomic.AddInt64(&env.allocBytesLeft, -size)
	if n < 0 {
		// TODO(quasilyte): provide JVM stack trace info.
		panic("allocations limit reached")
	}
}
