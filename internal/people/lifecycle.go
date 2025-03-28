package people

import (
	"math/rand/v2"

	"github.com/janithl/citylyf/internal/entities"
)

func SimulateLifecycle() {
	for _, person := range entities.Sim.People.People {

		// --- Marriage ---
		if person.Relationship != entities.Married && person.Age() >= entities.AgeOfAdulthood {
			// Calculate the probability of marriage for the current person's age
			marriageProbability := marriageProbabilityByAge(float64(person.Age()), entities.ProbabilityOfMarriage/entities.DaysPerYear)
			if rand.Float64() < marriageProbability {
				if candidate := findMarriageCandidate(person); candidate != nil {
					Marry(person, candidate)
				}
			}
		}
	}
}
