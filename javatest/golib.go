package javatest

import (
	"bytes"
	"fmt"
)

var golibOutput bytes.Buffer

func golibPrintInt(x int32) {
	fmt.Fprintf(&golibOutput, "%d\n", x)
}

func golibPrintLong(x int64) {
	fmt.Fprintf(&golibOutput, "%d\n", x)
}
