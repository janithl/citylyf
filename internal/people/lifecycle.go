package people

import (
	"fmt"
	"math/rand/v2"

	"github.com/janithl/citylyf/internal/entities"
)

func SimulateLifecycle() {
	for _, person := range entities.Sim.People.People {
		// --- Retirement ---
		// Assume a normal distribution for the age of retirement
		retirementAge := entities.MeanRetirementAge + rand.NormFloat64()*entities.StdDevRetirementAge
		if person.Age() >= int(retirementAge) && person.CareerLevel != entities.Retired &&
			rand.Float64() < 1/(entities.DaysPerYear*entities.StdDevRetirementAge*2) { // probability of retirement is spread out over a 5 year period
			if company, ok := entities.Sim.Companies[person.EmployerID]; ok {
				company.RemoveEmployee(person.ID)
			}
			person.EmployerID = 0
			person.CareerLevel = entities.Retired
			fmt.Printf("[  Job ] %s %s (%d) has retired\n", person.FirstName, person.FamilyName, person.Age())
		}

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
