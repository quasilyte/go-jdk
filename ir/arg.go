package ir

import (
	"fmt"
	"math"
)

// Arg is an instruction argument (sometimes called operand).
// Value bits (int64) should be interpreted differently depending on the Kind.
type Arg struct {
	Kind  ArgKind
	Value int64
}

func (arg Arg) FloatValue() float32 { return math.Float32frombits(uint32(arg.Value)) }

func (arg Arg) DoubleValue() float64 { return math.Float64frombits(uint64(arg.Value)) }

func (arg Arg) String() string {
	switch arg.Kind {
	case ArgBranch:
		return fmt.Sprintf("@%d", arg.Value)
	case ArgFlags:
		return "flags"
	case ArgReg:
		return fmt.Sprintf("r%d", arg.Value)
	case ArgIntConst, ArgLongConst:
		return fmt.Sprintf("%d", arg.Value)
	case ArgFloatConst:
		return formatFloat64(float64(arg.FloatValue()))
	case ArgDoubleConst:
		return formatFloat64(arg.DoubleValue())
	default:
		return fmt.Sprintf("{%d,%d}", arg.Kind, arg.Value)
	}
}

// ArgKind describes an argument category.
// Most kinds change how Arg value should be interpreted.
type ArgKind int

const (
	ArgInvalid ArgKind = iota

	ArgBranch
	ArgFlags
	ArgReg
	ArgIntConst
	ArgLongConst
	ArgFloatConst
	ArgDoubleConst
)
