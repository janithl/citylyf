package people

import (
	"math"
	"math/rand"
	"time"
)

// getAge generates a random age based on a bell curve
func getAge(mean, stdDev float64, minAge, maxAge int) int {
	if minAge < 0 {
		minAge = 0
	}
	if maxAge < minAge {
		maxAge = minAge
	}

	maxCalculations := 100
	for i := 0; i < maxCalculations; i++ {
		// Box-Muller transform to generate normal distribution
		u1 := rand.Float64()
		u2 := rand.Float64()
		z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

		// Scale and shift to get the desired mean and standard deviation
		age := mean + z*stdDev

		// Ensure age is within bounds
		if int(age) >= minAge && int(age) <= maxAge {
			return int(math.Round(age))
		}
	}

	return int(minAge + rand.Intn(maxAge))
}

// getRandomBirthdate generates a random birthdate given the age
func getRandomBirthdate(age int) time.Time {
	year := time.Now().Year() - age

	// Generate a random month (1-12)
	month := time.Month(rand.Intn(12) + 1)

	// Generate a random day based on the month and year
	// Use time.Date to determine the last day of the month
	day := rand.Intn(time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()) + 1

	birthdate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if time.Now().Before(birthdate) {
		return time.Now()
	} else {
		return birthdate
	}
}
