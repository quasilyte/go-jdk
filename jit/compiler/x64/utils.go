package x64

import (
	"github.com/quasilyte/go-jdk/ir"
)

func regDisp(arg ir.Arg) int32 {
	return int32(arg.Value * 8)
}
