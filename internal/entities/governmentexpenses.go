package entities

import (
	"math"

	"github.com/janithl/citylyf/internal/utils"
)

type CostType string

const (
	AsphaltRoadConstruction  CostType = "AsphaltRoadConstruction"
	AsphaltRoadMaintenance   CostType = "AsphaltRoadMaintenance"
	UnsealedRoadConstruction CostType = "UnsealedRoadConstruction"
	UnsealedRoadMaintenance  CostType = "UnsealedRoadMaintenance"
)

// GetGovernmentSpending returns government capex + opex spending in millions of dollars
func (g *Government) GetGovernmentSpending() float64 {
	return float64(utils.GetLastValue(g.CapExValues)+utils.GetLastValue(g.OpExValues)) / 1e6
}

// GetCapEx gets a particular capital expense
func (g *Government) GetCapEx(costType CostType, units int) int {
	if unitCost, exists := g.Expenses[costType]; exists {
		return int(math.Ceil(unitCost * float64(units))) // round cost and add
	}
	return 0
}

// AddCapEx adds a particular capital expense
func (g *Government) AddCapEx(costType CostType, units int) {
	g.CapEx += g.GetCapEx(costType, units)
}

// calculate annual government opex
func (g *Government) CalculateOpEx() int {
	roadMaintenanceCost := 0
	for _, r := range Sim.Geography.GetRoads() {
		if r.Type == Asphalt {
			roadMaintenanceCost += r.GetLength() * int(g.Expenses[AsphaltRoadMaintenance])
		} else {
			roadMaintenanceCost += r.GetLength() * int(g.Expenses[UnsealedRoadMaintenance])
		}
	}

	return roadMaintenanceCost // TODO: Add other maintenance costs
}

// run annually to update expenses
func (g *Government) ReviseExpenses() {
	for costType, unitCost := range g.Expenses {
		g.Expenses[costType] += unitCost * Sim.Market.InflationRate() / 100
	}
}

func NewExpenses() map[CostType]float64 {
	expenses := make(map[CostType]float64)
	expenses[AsphaltRoadConstruction] = 15000
	expenses[AsphaltRoadMaintenance] = 225
	expenses[UnsealedRoadConstruction] = 135
	expenses[UnsealedRoadMaintenance] = 121

	return expenses
}
