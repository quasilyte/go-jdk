package x64

import (
	"math"

	"github.com/quasilyte/go-jdk/ir"
)

func ptrDisp(arg ir.Arg) int32 {
	return int32(arg.Value*16) + 8
}

func scalarDisp(arg ir.Arg) int32 {
	return int32(arg.Value * 16)
}

func regDisp(arg ir.Arg) int32 {
	return int32(arg.Value * 16)
}

func fits32bit(x int64) bool {
	return x >= math.MinInt32 && x <= math.MaxInt32
}
