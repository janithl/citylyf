package entities

import (
	"math"

	"github.com/janithl/citylyf/internal/utils"
)

type CostType string

const (
	AsphaltRoadConstruction CostType = "AsphaltRoadConstruction"
	AsphaltRoadMaintenance  CostType = "AsphaltRoadMaintenance"
)

// GetGovernmentSpending returns government capex + opex spending in millions of dollars
func (g *Government) GetGovernmentSpending() float64 {
	return float64(utils.GetLastValue(g.CapExValues)+utils.GetLastValue(g.OpExValues)) / 1e6
}

// AddCapEx adds a particular capital expense
func (g *Government) AddCapEx(costType CostType, units int) {
	if unitCost, exists := g.Expenses[costType]; exists {
		g.CapEx += int(math.Ceil(unitCost * float64(units))) // round cost and add
	}
}

// calculate annual government opex
func (g *Government) CalculateOpEx() int {
	roadMaintenanceUnitCost := g.Expenses[AsphaltRoadMaintenance]
	roadMaintenanceCost := 0
	for _, r := range Sim.Geography.GetRoads() {
		roadMaintenanceCost += r.GetLength() * int(roadMaintenanceUnitCost)
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

	return expenses
}
