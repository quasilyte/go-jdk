package irgen

import (
	"strings"

	"github.com/quasilyte/go-jdk/bytecode"
	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/jclass"
)

// isJump reports whether inst is a jump instruction.
func isJump(inst ir.Inst) bool {
	switch inst.Kind {
	case ir.InstJump, ir.InstJumpEqual, ir.InstJumpNotEqual, ir.InstJumpGtEq, ir.InstJumpGt, ir.InstJumpLt, ir.InstJumpLtEq:
		return true
	default:
		return false
	}
}

func isUnconditionalBranch(op bytecode.Op) bool {
	return op == bytecode.Goto
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

func findFrame(offset int, frames []jclass.StackMapFrame) *jclass.StackMapFrame {
	for i, frame := range frames {
		if int(frame.Offset) == offset {
			return &frames[i]
		}
	}
	return nil
}
