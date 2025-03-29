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
			entities.Sim.Companies.RemoveEmployeeFromTheirCompany(person)
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

		// --- Moving Out of Home ---
		if person.Age() >= entities.AgeOfAdulthood && person.Relationship != entities.Married {
			oldHousehold := entities.Sim.People.GetHouseholdByPersonID(person.ID)
			if oldHousehold == nil || oldHousehold.GetAdultCount() < 2 {
				continue // already head adult, no need to move out
			}

			if rand.Float64() < entities.ProbabilityOfMovingOut/entities.DaysPerYear { // annual rate spread out over each day of the year
				newHousehold := &entities.Household{
					ID:         entities.Sim.GetNextID(),
					MemberIDs:  []int{person.ID},
					MoveInDate: entities.Sim.Date,
					LastPayDay: entities.Sim.Date,
					Savings:    person.Savings,
				}
				if houseID := newHousehold.FindHousing(); houseID > 0 { // only move out if we can find new housing
					oldHousehold.RemoveMember(person)
					entities.Sim.People.Households[newHousehold.ID] = newHousehold
					fmt.Printf("[ Move ] %s %s (%d) has moved into house #%d, %d houses remain\n", person.FirstName,
						person.FamilyName, person.Age(), houseID, entities.Sim.Houses.GetFreeHouses())
				}
			}
		}
	}
}
