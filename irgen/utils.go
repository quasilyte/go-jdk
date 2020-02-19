package irgen

import (
	"strings"

	"github.com/quasilyte/GopherJRE/ir"
	"github.com/quasilyte/GopherJRE/jclass"
)

// isJump reports whether inst is a jump instruction.
func isJump(inst ir.Inst) bool {
	switch inst.Kind {
	case ir.InstJump, ir.InstJumpEqual, ir.InstJumpNotEqual, ir.InstJumpGtEq:
		return true
	default:
		return false
	}
}

func argsCount(d string) int {
	n := 0
	jclass.MethodDescriptor(d).WalkParams(func(jclass.DescriptorType) {
		n++
	})
	return n
}

func splitName(full string) (name, pkg string) {
	delim := strings.LastIndexByte(full, '/')
	if delim == -1 {
		return full, ""
	}
	return full[delim+1:], full[:delim]
}
