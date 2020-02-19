package ir

import (
	"strings"
)

// Inst is a single IR instruction.
//
// If instruction has explicit output, Dst represents
// output assignment destination.
type Inst struct {
	Kind  InstKind
	Dst   Arg
	Args  []Arg
	Flags InstFlags
}

func (inst Inst) String() string {
	var buf strings.Builder
	if inst.Dst.Kind != 0 {
		buf.WriteString(inst.Dst.String())
		buf.WriteString(" = ")
	}
	buf.WriteString(inst.Kind.String())
	for _, arg := range inst.Args {
		buf.WriteByte(' ')
		buf.WriteString(arg.String())
	}
	return buf.String()
}

// InstKind describes instruction operation.
// It can be treated as an "opcode" of IR abstract machine.
type InstKind int

//go:generate stringer -type=InstKind -trimprefix=Inst
const (
	InstInvalid InstKind = iota

	InstIload
	InstRet
	InstIret
	InstLret
	InstCallStatic
	InstIcmp
	InstLcmp
	InstJump
	InstJumpEqual
	InstJumpNotEqual
	InstJumpGtEq
	InstIadd
	InstLadd
	InstFadd
	InstIsub
	InstIneg
	InstLneg
	InstDadd
	InstConvL2I
	InstConvF2I
	InstConvD2I
)
