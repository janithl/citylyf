package entities

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/janithl/citylyf/internal/utils"
)

type People struct {
	Population       int
	PopulationValues []int // Historical values
	LabourForce      int   // Employable people
	Unemployed       int
	Households       []Household
}

func (p *People) UnemploymentRate() float64 {
	return 100.0 * float64(p.Unemployed) / float64(p.LabourForce)
}

func (p *People) PopulationGrowthRate() float64 {
	lastPopulationValue := p.PopulationValues[len(p.PopulationValues)-1]
	return 100.0 * float64(p.Population-lastPopulationValue) / float64(lastPopulationValue)
}

func (p *People) MoveIn(createHousehold func() Household) {
	for i := 0; i < rand.Intn(1+(Sim.Houses.GetFreeHouses()/4)); i++ {
		h := createHousehold()
		monthlyRentBudget := float64(h.AnnualIncome()) / (4 * 12)              // 25% of yearly income towards rent / 12
		houseId := Sim.Houses.MoveIn(int(monthlyRentBudget), len(h.Members)/2) // everyone gets to share a bedroom
		if houseId > 0 {
			h.HouseID = houseId
			fmt.Printf("[ Move ] %s family has moved into a house, %d houses remain\n", h.FamilyName(), Sim.Houses.GetFreeHouses())
			p.Households = append(p.Households, h)
			p.Population += len(h.Members)
		}
	}
}

func (p *People) MoveOut() {
	h := Sim.People.Households
	// traverse in reverse order to avoid index shifting
	for i := len(h) - 1; i >= 0; i-- {
		if len(h[i].Members) > 0 && h[i].IsEligibleForMoveOut() {
			movedName := h[i].FamilyName()
			h = slices.Delete(h, i, i+1)
			Sim.Houses.MoveOut()
			fmt.Printf("[ Move ] %s family has moved out of the city, %d houses remain\n", movedName, Sim.Houses.GetFreeHouses())
		}
	}
	Sim.People.Households = h
}

// calculate the unemployed and the total labour force
func (p *People) CalculateUnemployment() {
	labourforce, unemployed := 0, 0
	for i := 0; i < len(p.Households); i++ {
		for j := 0; j < len(p.Households[i].Members); j++ {
			if p.Households[i].Members[j].IsEmployable() {
				labourforce += 1
				if !p.Households[i].Members[j].IsEmployed() {
					unemployed += 1
				}
			}
		}
	}
	p.LabourForce = labourforce
	p.Unemployed = unemployed
}

// Append current population value to history
func (p *People) UpdatePopulationValues() {
	p.PopulationValues = utils.AddFifo(p.PopulationValues, p.Population, 20)
}
