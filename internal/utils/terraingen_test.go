package utils

import (
	"testing"
)

// TestGenerateElevationMap checks that Generate produces a valid elevation map
func TestGenerateElevationMap(t *testing.T) {
	size := 16
	max := 10
	sea := 5

	elevationMap := GenerateElevationMap(sea, max, size, 0.0015, 0.005, 0.01)
	if len(elevationMap) != size {
		t.Errorf("Map size is %d, expected %d", len(elevationMap), size)
	}

	for x := range elevationMap {
		for y := range elevationMap[x] {
			if elevationMap[x][y] > max {
				t.Errorf("Map elevation at (%d, %d) is %d, expected maximum is %d", x, y, elevationMap[x][y], max)
			}
			if elevationMap[x][y] < 0 {
				t.Errorf("Map elevation at (%d, %d) is %d, expected minimum is 0", x, y, elevationMap[x][y])
			}
		}
	}
}
