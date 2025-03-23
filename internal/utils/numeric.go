package utils

import (
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
func GetLastValue[T int | float64](s []T) T {
	lastIndex := len(s) - 1
	if lastIndex >= 0 {
		return s[lastIndex]
	}
	return 0
}

// Adds an element FIFO for float64 and int
func AddFifo[T int | float64 | string](s []T, element T, maxLength int) []T {
	if len(s) >= maxLength {
		s = slices.Delete(s, 0, 1) // Remove oldest element
	}
	return append(s, element)
}

// Check if a target is within a range
func IsWithinRange(limit1, limit2, target int) bool {
	return min(limit1, limit2) <= target && target <= max(limit1, limit2)
}
