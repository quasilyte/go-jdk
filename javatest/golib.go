package javatest

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/quasilyte/go-jdk/jruntime"
)

var golibOutput bytes.Buffer

func golibNop() {}

func golibPrintInt(x int32) {
	fmt.Fprintf(&golibOutput, "%d\n", x)
}

func golibPrintLong(x int64) {
	fmt.Fprintf(&golibOutput, "%d\n", x)
}

func golibPrintIntArray(xs *jruntime.IntArrayObject) {
	var parts []string
	for _, x := range xs.AsSlice() {
		p := strconv.Itoa(int(x))
		parts = append(parts, p)
	}
	fmt.Fprintf(&golibOutput, "[%s]\n", strings.Join(parts, ", "))
}

func golibIsub(x, y int32) int32 {
	return x - y
}

func golibIsub3(x, y, z int32) int32 {
	return x - y - z
}

func golibII_L(a1 int32, a2 int32) int64 {
	return int64(a1 - a2)
}

func golibLI_I(a1 int64, a2 int32) int32 {
	return int32(a1) - a2
}

func golibIL_I(a1 int32, a2 int64) int32 {
	return a1 - int32(a2)
}

func golibILIL_I(a1 int32, a2 int64, a3 int32, a4 int64) int32 {
	return a1 - int32(a2) - a3 - int32(a4)
}

func golibGC() {
	runtime.GC()
}
