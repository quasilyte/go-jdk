package x64

import (
	"fmt"
	"testing"
)

func TestLongJumps(t *testing.T) {
	checkEncoding := func(t *testing.T, result []byte, want string) {
		have := fmt.Sprintf("%x", result)[:len(want)]
		if have != want {
			t.Errorf("encoding mismatch:\nhave: %s\nwant: %s", have, want)
		}
	}

	t.Run("jmp", func(t *testing.T) {
		asm := NewAssembler()
		asm.Jmp(0)
		asm.Nop(0x9988)
		asm.Label(0)
		checkEncoding(t, asm.Link(), "e988990000")
	})

	t.Run("jge", func(t *testing.T) {
		asm := NewAssembler()
		asm.Jge(0)
		asm.Nop(0x6677)
		asm.Label(0)
		checkEncoding(t, asm.Link(), "0f8d77660000")
	})
}
