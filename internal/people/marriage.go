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

	p1household := entities.Sim.People.GetHouseholdByPersonID(person1.ID)
	p2household := entities.Sim.People.GetHouseholdByPersonID(person2.ID)

	if p1household != nil {
		if p1household.GetAdultCount() == 1 {
			// Person1 is the only adult in the household, so add Person2 to the same household
			p1household.AddMember(person2.ID, person2.Savings)
			fmt.Printf("[ Weds ] %s moves in with the %s family\n", person2.FirstName, p1household.FamilyName())
		} else {
			p1household.RemoveMember(person1)
		}
	}

	if p2household != nil {
		if p2household.GetAdultCount() == 1 {
			// Person2 is the only adult in the household, so add Person1 to the same household or combine households
			if p1household.IsMember(person2.ID) {
				for _, id := range p2household.MemberIDs {
					p1household.AddMember(id, 0)
				}
				fmt.Printf("[ Weds ] %s and %s families combine\n", p1household.FamilyName(), p2household.FamilyName())
				delete(entities.Sim.People.Households, p2household.ID)
			} else {
				p2household.AddMember(person1.ID, person1.Savings)
				fmt.Printf("[ Weds ] %s moves in with the %s family\n", person1.FirstName, p2household.FamilyName())
			}
		} else {
			p2household.RemoveMember(person2)
		}
	}

	// if we already have a household, we can stop
	// else we have to create a new household for the couple
	if (p1household != nil && p1household.IsMember(person2.ID)) || (p2household != nil && p2household.IsMember(person1.ID)) {
		return
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

// marriageProbabilityByAge calculates marriage probability based on age using a normal distribution
func marriageProbabilityByAge(currentAge, maxProb float64) float64 {
	// We'll use the probability density function (PDF) of the normal distribution
	exponent := -0.5 * math.Pow((currentAge-entities.MeanMarriageAge)/entities.StdDevMarriageAge, 2)
	coefficient := 1 / (entities.StdDevMarriageAge * math.Sqrt(2*math.Pi))
	pdf := coefficient * math.Exp(exponent)

	// Scale the PDF to get a probability, peaking at maxProb around the mean
	// We can normalize it by the PDF at the mean age to ensure the peak is at maxProb
	pdfAtMean := 1 / (entities.StdDevMarriageAge * math.Sqrt(2*math.Pi))
	if pdfAtMean > 0 {
		probability := (pdf / pdfAtMean) * maxProb
		return math.Max(0, math.Min(1, probability)) // Ensure probability is between 0 and 1
	}
	return 0
}
