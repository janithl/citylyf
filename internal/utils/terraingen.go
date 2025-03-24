package utils

import (
	"errors"
	"math"
	"math/rand/v2"
)

// GetElevationSlice returns a slice of elevations starting at max and decaying
// exponentially toward min. The decayFactor should be between 0 and 1, where values
// closer to 1 result in slower decay.
func GetElevationSlice(max, min, maxNum int, decayFactor float64) []int {
	if max <= min {
		return []int{min}
	}

	var elevations []int
	current := float64(max)
	rounded := int(math.Round(current))
	minimum := float64(min)
	decayFactor = math.Min(decayFactor, 0.99) // decay factor cannot be more than 0.99

	// Append the starting elevation.
	elevations = append(elevations, int(math.Round(current)))

	// Generate intermediate levels until we get close enough to min.
	// We stop when the next computed level would be the same as min (or lower).
	for rounded > min && len(elevations) < maxNum {
		// Exponential decay: move current partway toward min.
		current = minimum + (current-minimum)*decayFactor
		rounded = int(math.Round(current))
		elevations = append(elevations, rounded)
	}

	// Ensure min is included.
	if elevations[len(elevations)-1] != min {
		elevations = append(elevations, min)
	}

	return elevations
}

// GenerateElevationMap generates the elevation values on the map
// From: https://janithl.github.io/2019/09/go-terrain-gen-part-4/
func GenerateElevationMap(seaLevel, maxElevation, size int, peakProbability, rangeProbability, cliffProbability float64) [][]int {
	elevationSteps := GetElevationSlice(maxElevation, seaLevel-1, 20, 0.7+peakProbability*50)
	elevationSteps = append(elevationSteps, GetElevationSlice(seaLevel-1, 0, 10, 0.3+peakProbability*100)...)

	elevations := make([][]int, size)
	for i := range size {
		elevations[i] = make([]int, size)
	}

	// bias x and y create a vector along which mountain ranges form
	biasX := rand.IntN(6) - 3
	biasY := rand.IntN(6) - 3

	// iterate down from max elevation, assigning vals
	for _, e := range elevationSteps {
		for x := range size {
			for y := range size {
				// if the element is next to a element with elevation x, it
				// should get elevation x - 1
				// alternately, if the random value meets our criteria, it's a peak
				if GetAdjacentElevation(elevations, x, y, e, cliffProbability) || rand.Float64() < peakProbability {
					setElevation(elevations, x, y, e)
					if rand.Float64() > rangeProbability { // randomly add follow-up peaks
						setElevation(elevations, x+biasX, y+biasY, e)
					}
					if rand.Float64() > rangeProbability {
						setElevation(elevations, x-biasX, y-biasY, e)
					}
				}
			}
		}
	}

	return elevations
}

// adjacentElevation checks if an adjacent element
// to the given element (h, w) is at a given elevation
func GetAdjacentElevation(elevations [][]int, w, h, elevation int, cliffProbability float64) bool {
	for x := w - 1; x <= w+1; x++ {
		for y := h - 1; y <= h+1; y++ {
			if x == w && y == h {
				continue
			}

			if currentElevation, err := getElevation(elevations, x, y); err == nil && currentElevation == elevation+1 {
				// if this element is *not* randomly a cliff, return true
				return rand.Float64() > cliffProbability
			}
		}
	}

	return false
}

func setElevation(elevations [][]int, x, y, e int) {
	if currentElevation, err := getElevation(elevations, x, y); err == nil && currentElevation == 0 {
		elevations[x][y] = e
	}
}

func getElevation(elevations [][]int, x, y int) (int, error) {
	if x < 0 || y < 0 || x >= len(elevations) || y >= len(elevations) {
		return 0, errors.New("index out of bounds")
	}
	return elevations[x][y], nil
}
