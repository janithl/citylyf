package utils

import (
	"fmt"
	"math"
	"slices"
)

// Convert a slice of int to float64
func ConvertToF64(s []int) []float64 {
	f64 := make([]float64, len(s))
	for i, val := range s {
		f64[i] = float64(val)
	}
	return f64
}

// Get the last value of a float64 slice
func GetLastValue(s []float64) float64 {
	lastIndex := len(s) - 1
	if lastIndex >= 0 {
		return s[lastIndex]
	}
	return 0
}

// Adds an element FIFO
func AddFifo(s []float64, element float64, maxLength int) []float64 {
	if len(s) >= maxLength {
		s = slices.Delete(s, 0, 1) // FIFO behavior (oldest values removed)
	}
	return append(s, element)
}

// Format currency with suffixes
func FormatCurrency(value float64, currencySymbol string) string {
	suffix := " "
	if math.Abs(value) > 1e6 {
		value /= 1e6
		suffix = "M"
	} else if math.Abs(value) > 1e3 {
		value /= 1e3
		suffix = "K"
	}

	return fmt.Sprintf("%s %.2f %s", currencySymbol, value, suffix)
}
