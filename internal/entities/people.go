package entities

import (
	"maps"
	"math"
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

// GetHouseholdByPersonID returns the household a person belongs to
func (p *People) GetHouseholdByPersonID(personID int) *Household {
	for _, household := range p.Households {
		if household.IsMember(personID) {
			return household
		}
	}
	return nil
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

func (p *People) AverageMonthlyDisposableIncome() int {
	if len(p.Households) == 0 {
		return 0 // Avoid division by zero
	}

	totalDisposableIncome := 0.0
	for _, household := range p.Households {
		disposable := float64(household.AnnualIncome(false))/12.0 - float64(household.LastMonthExpenses)
		if disposable < 0 {
			disposable = 0
		}
		totalDisposableIncome += disposable
	}

	return int(math.Round(totalDisposableIncome / float64(len(p.Households))))
}
