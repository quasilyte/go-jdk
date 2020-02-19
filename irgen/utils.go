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

func argsCount(d string) int {
	d = d[len("("):]
	n := 0
	for d[0] != ')' {
		skip := 1
		switch d[0] {
		case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z':
			n++
		case 'L':
			skip = strings.IndexByte(d, ';')
			n++
		case '[':
			// Do nothing.
		}
		d = d[skip:]
	}
	return n
}

func splitName(full string) (name, pkg string) {
	delim := strings.LastIndexByte(full, '/')
	if delim == -1 {
		return full, ""
	}
	return full[delim+1:], full[:delim]
}
