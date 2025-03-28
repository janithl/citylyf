package people

import (
	"math/rand/v2"

	"github.com/janithl/citylyf/internal/entities"
)

func SimulateLifecycle() {
	for _, person := range entities.Sim.People.People {
		// --- Marriage ---
		if person.Relationship != entities.Married && person.Age() >= entities.AgeOfAdulthood {
			// Probability of marriage can also be age-dependent, but for simplicity,
			// we'll use a fixed probability here. You could model this with a distribution too.
			if rand.Float64() < entities.ProbabilityOfMarriage/entities.DaysPerYear {
				if candidate := findMarriageCandidate(person); candidate != nil {
					Marry(person, candidate)
				}
			}
		}
	}
}
