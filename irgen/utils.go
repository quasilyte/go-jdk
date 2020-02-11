package irgen

import (
	"strings"

	"github.com/quasilyte/GopherJRE/ir"
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

func splitName(full string) (name, pkg string) {
	delim := strings.LastIndexByte(full, '/')
	if delim == -1 {
		return full, ""
	}
	return full[delim+1:], full[:delim]
}
