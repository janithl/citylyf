package utils

import (
	"fmt"
	"math"
)

// Format currency with suffixes
func FormatCurrency(value float64, currencySymbol string) string {
	suffix := " "
	if math.Abs(value) > 1e9 {
		value /= 1e9
		suffix = "B"
	} else if math.Abs(value) > 1e6 {
		value /= 1e6
		suffix = "M"
	} else if math.Abs(value) > 1e3 {
		value /= 1e3
		suffix = "K"
	}

	return fmt.Sprintf("%s %.2f %s", currencySymbol, value, suffix)
}
