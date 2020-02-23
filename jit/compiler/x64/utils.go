package x64

import (
	"math"

	"github.com/quasilyte/go-jdk/ir"
)

func regDisp(arg ir.Arg) int32 {
	return int32(arg.Value * 8)
}

func fits32bit(x int64) bool {
	return x >= math.MinInt32 && x <= math.MaxInt32
}
