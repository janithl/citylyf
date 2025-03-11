package people

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/janithl/citylyf/internal/entities"
)

// getAge generates a random age based on a bell curve
func getAge(mean, stdDev float64, minAge, maxAge int) (years int, months int) {
	if minAge < 0 {
		minAge = 0
	}
	if maxAge-minAge < 1 {
		return minAge, 0
	}

	maxCalculations := 100
	for i := 0; i < maxCalculations; i++ {
		age := mean + rand.NormFloat64()*stdDev

		// Ensure age is within bounds
		if age >= float64(minAge) && age <= float64(maxAge) {
			return int(math.Round(age)), 1 + (int(age*12) % 12)
		}
	}

	return int(minAge + rand.IntN(maxAge-minAge)), 1 + rand.IntN(12)
}

// getRandomBirthdate generates a random birthdate given the age
func getRandomBirthdate(ageY int, ageM int) time.Time {
	currentDate := entities.Sim.Date
	year := currentDate.Year() - ageY

	// Generate a random day based on the month and year
	// Use time.Date to determine the last day of the month
	day := rand.IntN(time.Date(year, time.Month(ageM), 0, 0, 0, 0, 0, time.UTC).Day()) + 1
	birthdate := time.Date(year, time.Month(ageM), day, 0, 0, 0, 0, time.UTC)

	if currentDate.Before(birthdate) {
		return currentDate.AddDate(0, -rand.IntN(12), -rand.IntN(28))
	} else {
		return birthdate
	}
}
