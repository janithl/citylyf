package utils

import (
	"fmt"
	"math"
)

func GetTurningPoint(x1, y1, x2, y2 int) (int, int) {
	xdiff := x1 - x2
	ydiff := y1 - y2

	if xdiff < 0 {
		xdiff = -xdiff
	}
	if ydiff < 0 {
		ydiff = -ydiff
	}

	if xdiff < ydiff {
		return x1, y2
	}

	return x2, y1
}

// FormatDistance
func FormatDistance(value float64) string {
	suffix := "m"
	if math.Abs(value) > 1e3 {
		value /= 1e3
		suffix = "km"
	}

	return fmt.Sprintf("%.1f %s", value, suffix)
}
