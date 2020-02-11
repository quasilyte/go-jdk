package ir

import (
	"strconv"
)

// formatFloat64 returns string representation of f that can't be
// mixed up with integer literal.
func formatFloat64(f float64) string {
	switch f {
	case 0:
		return "0.0"
	case 1:
		return "1.0"
	default:
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
}
