package entities

import (
	"fmt"
	"maps"
	"math/rand"
	"slices"

	"github.com/janithl/citylyf/internal/utils"
)

const AgeGroupSize = 10

type AgeGroup struct {
	Male   int
	Female int
	Other  int
}

type People struct {
	PopulationValues       []int // Historical values
	LabourForce            int   // Employable people
	Unemployed             int
	UnemploymentRateValues []float64
	People                 map[int]*Person
	Households             map[int]*Household
	AgeGroups              map[int]AgeGroup // Population breakdown by age group
}

func (p *People) Population() int {
	return len(p.People)
}

func (p *People) UnemploymentRate() float64 {
	if p.LabourForce == 0 {
		return 0.0
	}
	return 100.0 * float64(p.Unemployed) / float64(p.LabourForce)
}

func (p *People) PopulationGrowthRate() float64 {
	lastPopulationValue := p.PopulationValues[len(p.PopulationValues)-1]
	if lastPopulationValue == 0 {
		return 0.0
	}
	return 100.0 * float64(p.Population()-lastPopulationValue) / float64(lastPopulationValue)
}

func (p *People) MoveIn(createHousehold func() *Household) {
	if Sim.Houses.GetFreeHouses() == 0 || rand.Float64() < 0.85 { // 15% change of moving in if there are free houses
		return
	}

	h := createHousehold()
	monthlyRentBudget := float64(h.AnnualIncome()) / (4 * 12)        // 25% of yearly income towards rent / 12
	houseId := Sim.Houses.MoveIn(int(monthlyRentBudget), h.Size()/2) // everyone gets to share a bedroom
	if houseId > 0 {
		h.HouseID = houseId
		fmt.Printf("[ Move ] %s family has moved into a house, %d houses remain\n", h.FamilyName(), Sim.Houses.GetFreeHouses())
		p.Households[h.ID] = h
	}
}

// GetPerson gets an existing person
func (p *People) GetPerson(personID int) *Person {
	person, exists := p.People[personID]
	if exists {
		return person
	}
	return nil
}

// AddPerson adds a new person
func (p *People) AddPerson(person *Person) {
	p.People[person.ID] = person
}

// RemovePerson removes person
func (p *People) RemovePerson(personID int) {
	delete(p.People, personID)
}

// GetHouseholdIDs returns a sorted list of household IDs
func (p *People) GetHouseholdIDs() []int {
	IDs := []int{}
	for household := range maps.Values(p.Households) {
		IDs = append(IDs, household.ID)
	}
	slices.Sort(IDs)
	return IDs
}

func (p *People) MoveOut(removeEmployeeFromCompany func(companyID int, employeeID int)) {
	for household := range maps.Values(p.Households) {
		if household.Size() > 0 && household.IsEligibleForMoveOut() {
			movedName := household.FamilyName()
			// remove members from people, their jobs, and deduct from population
			for _, memberID := range household.MemberIDs {
				member := p.GetPerson(memberID)
				if member != nil {
					removeEmployeeFromCompany(member.EmployerID, memberID)
					p.RemovePerson(memberID)
				}
			}
			delete(p.Households, household.ID)
			Sim.Houses.MoveOut()
			fmt.Printf("[ Move ] %s family has moved out of the city, %d houses remain\n", movedName, Sim.Houses.GetFreeHouses())
		}
	}
}

// calculate the unemployed and the total labour force
func (p *People) CalculateUnemployment() {
	labourforce, unemployed := 0, 0
	for _, person := range p.People {
		if person.IsEmployable() {
			labourforce += 1
			if !person.IsEmployed() {
				unemployed += 1
			}
		}
	}
	p.LabourForce = labourforce
	p.Unemployed = unemployed
	p.UnemploymentRateValues = utils.AddFifo(p.UnemploymentRateValues, p.UnemploymentRate(), 20)
}

// Append current population value to history
func (p *People) UpdatePopulationValues() {
	p.PopulationValues = utils.AddFifo(p.PopulationValues, p.Population(), 20)
}

// calculate the age groups of the population
func (p *People) CalculateAgeGroups() {
	groups := make(map[int]AgeGroup)
	for i := 0; i < 120; i += AgeGroupSize {
		groups[i] = AgeGroup{}
	}

	for _, person := range p.People {
		ageGroup := AgeGroupSize * (person.Age() / AgeGroupSize)
		if group, ok := groups[ageGroup]; ok {
			switch person.Gender {
			case Male:
				group.Male += 1
			case Female:
				group.Female += 1
			case Other:
				group.Other += 1
			}
			groups[ageGroup] = group
		}
	}
	p.AgeGroups = groups
}
