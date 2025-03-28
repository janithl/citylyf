package people

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/janithl/citylyf/internal/entities"
)

// Marry simulates marriage between two people.
func Marry(person1 *entities.Person, person2 *entities.Person) {
	if person1.Relationship == entities.Married || person2.Relationship == entities.Married ||
		person1.Age() < entities.AgeOfAdulthood || person2.Age() < entities.AgeOfAdulthood {
		return // cannot marry if already married or not an adult
	}

	person1.Relationship = entities.Married
	person2.Relationship = entities.Married

	fmt.Printf("[ Weds ] Wedding bells as %s %s (%d) marries %s %s (%d)!\n", person1.FirstName,
		person1.FamilyName, person1.Age(), person2.FirstName, person2.FamilyName, person2.Age())

	// TODO: Make sure no kids are left behind in these households, if the person is the only adult add their partner to the same household?
	if p1household := entities.Sim.People.GetHouseholdByPersonID(person1.ID); p1household != nil {
		p1household.RemoveMember(person1)
	}
	if p2household := entities.Sim.People.GetHouseholdByPersonID(person2.ID); p2household != nil {
		p2household.RemoveMember(person2)
	}

	household := &entities.Household{
		ID:         entities.Sim.GetNextID(),
		MemberIDs:  []int{person1.ID, person2.ID},
		MoveInDate: entities.Sim.Date,
		LastPayDay: entities.Sim.Date,
		Savings:    person1.Savings + person2.Savings,
	}

	if houseID := household.FindHousing(); houseID > 0 {
		entities.Sim.People.Households[household.ID] = household
	} else {
		fmt.Printf("[ Move ] The newlywed %s family has been unable to find housing, and has moved out of the city\n", household.FamilyName())
		RemoveHousehold(household)
	}
}

// findMarriageCandidates finds a suitable marriage candidate for a person.
func findMarriageCandidate(person *entities.Person) *entities.Person {
	eligibleCandidates := []*entities.Person{}
	for _, candidate := range entities.Sim.People.People {
		if candidate.ID != person.ID &&
			candidate.Relationship != entities.Married &&
			candidate.Age() > entities.AgeOfAdulthood &&
			candidate.FamilyName != person.FamilyName && // Sorry, George-Michael!
			math.Abs(float64(person.Age()-candidate.Age())) < entities.MaxMarriageAgeDifference { // Age difference within a reasonable range
			eligibleCandidates = append(eligibleCandidates, candidate)
		}
	}

	if len(eligibleCandidates) > 0 {
		return eligibleCandidates[rand.Intn(len(eligibleCandidates))]
	}

	return nil
}
