package utils

import (
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

// CalculateProbabilityByAge calculates a probability based on age using a normal distribution
func CalculateProbabilityByAge(mean, stdDev, currentAge, maxProb float64) float64 {
	// We'll use the probability density function (PDF) of the normal distribution
	exponent := -0.5 * math.Pow((currentAge-mean)/stdDev, 2)
	coefficient := 1 / (stdDev * math.Sqrt(2*math.Pi))
	pdf := coefficient * math.Exp(exponent)

	// Scale the PDF to get a probability, peaking at maxProb around the mean
	// We can normalize it by the PDF at the mean age to ensure the peak is at maxProb
	pdfAtMean := 1 / (stdDev * math.Sqrt(2*math.Pi))
	if pdfAtMean > 0 {
		probability := (pdf / pdfAtMean) * maxProb
		return math.Max(0, math.Min(1, probability)) // Ensure probability is between 0 and 1
	}
	return 0
}
