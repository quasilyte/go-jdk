package javatest

import (
	"bytes"
	"fmt"
)

var golibOutput bytes.Buffer

func golibNop() {}

func golibPrintInt(x int32) {
	fmt.Fprintf(&golibOutput, "%d\n", x)
}

func golibPrintLong(x int64) {
	fmt.Fprintf(&golibOutput, "%d\n", x)
}

func golibIsub(x, y int32) int32 {
	return x - y
}

func golibIsub3(x, y, z int32) int32 {
	return x - y - z
}
