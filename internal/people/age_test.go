package people

import (
	"math"
	"testing"
)

// TestGetAge checks that getAge produces values within the expected range.
func TestGetAge(t *testing.T) {
	mean := 30.0
	stdDev := 10.0
	minAge := 18
	maxAge := 65
	iterations := 10000

	for i := 0; i < iterations; i++ {
		ageY, ageM := getAge(mean, stdDev, minAge, maxAge)

		if ageY < minAge || ageY > maxAge {
			t.Errorf("getAge returned %dY %dM, which is outside the range [%d, %d]", ageY, ageM, minAge, maxAge)
		}
	}
}

// TestGetAgeEdgeCases checks handling of min/max edge cases.
func TestGetAgeEdgeCases(t *testing.T) {
	tests := []struct {
		mean, stdDev   float64
		minAge, maxAge int
	}{
		{30, 10, 20, 40}, // Normal range
		{30, 10, 30, 30}, // Single possible age
		{30, 10, 40, 20}, // Swapped min/max should return minAge (corrected)
		{30, 10, -5, 10}, // Negative minAge should be corrected to 0
	}

	for _, tt := range tests {
		ageY, ageM := getAge(tt.mean, tt.stdDev, tt.minAge, tt.maxAge)

		expectedMin := int(math.Max(float64(tt.minAge), 0)) // Ensure no negative min
		expectedMax := int(math.Max(float64(tt.maxAge), float64(expectedMin)))

		if ageY < expectedMin || ageY > expectedMax {
			t.Errorf("getAge(%v, %v, %v, %v) returned %dY %dM, expected range [%d, %d]",
				tt.mean, tt.stdDev, tt.minAge, tt.maxAge, ageY, ageM, expectedMin, expectedMax)
		}
	}
}
