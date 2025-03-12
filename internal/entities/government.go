package entities

import (
	"fmt"
	"maps"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

type Government struct {
	Reserves            int
	LastCalculationYear int
	CorporateTaxRate    float64      // Flat corporate tax rate
	IncomeTaxBrackets   []TaxBracket // Progressive income tax brackets

	// Historical values
	ReserveValues, CollectedTaxValues []int
}

// TaxBracket defines an income range and corresponding tax rate
type TaxBracket struct {
	Threshold int     // Income threshold for this bracket
	Rate      float64 // Tax rate for income above the threshold
}

// CollectTaxes runs annually, collecting from companies and households
func (g *Government) CollectTaxes() {
	if g.LastCalculationYear >= Sim.Date.Year() { // has already run this year
		return
	}

	// Collect household income taxes
	personalTaxesCollected := 0
	for household := range maps.Values(Sim.People.Households) {
		householdTax := g.CalculateIncomeTax(household.AnnualIncome(false))

		// Deduct tax from household wealth
		household.Savings -= householdTax
		personalTaxesCollected += householdTax
	}
	fmt.Printf("[  Tax ] Collected $%d in personal income taxes\n", personalTaxesCollected)

	// Collect corporate tax and reset tax payable account
	corporateTaxesCollected := 0
	for id := range Sim.Companies {
		corporateTaxesCollected += int(Sim.Companies[id].TaxPayable)
		Sim.Companies[id].TaxPayable = 0.0
	}
	fmt.Printf("[  Tax ] Collected $%d in corporate taxes\n", corporateTaxesCollected)

	// add collected taxes to government reserves
	totalTaxesCollected := personalTaxesCollected + corporateTaxesCollected
	g.Reserves += totalTaxesCollected
	fmt.Printf("[  Tax ] Collected $%d in total taxes for %d. Government reserves: $%d\n",
		totalTaxesCollected, g.LastCalculationYear, g.Reserves)
	g.LastCalculationYear = Sim.Date.Year()

	// Append current values to history
	g.ReserveValues = utils.AddFifo(g.ReserveValues, g.Reserves, 10)
	g.CollectedTaxValues = utils.AddFifo(g.CollectedTaxValues, totalTaxesCollected, 10)
}

// CalculateIncomeTax applies progressive tax rates to household income
func (g *Government) CalculateIncomeTax(income int) int {
	totalTax := 0.0
	for _, bracket := range g.IncomeTaxBrackets {
		if income > bracket.Threshold {
			taxableAmount := float64(income - bracket.Threshold)
			if taxableAmount > 0 {
				totalTax += taxableAmount * (bracket.Rate / 100)
				income = bracket.Threshold
			}
		}
	}
	return int(totalTax)
}

// TODO: Replace with actual spending calculation
func (g *Government) GetGovernmentSpending() float64 {
	return 5.0
}

// NewGovernment initializes the government system with reserves and progressive tax brackets
func NewGovernment(reserves int, startDate time.Time) *Government {
	return &Government{
		Reserves:            reserves,
		LastCalculationYear: startDate.Year(),
		CorporateTaxRate:    9.5,
		IncomeTaxBrackets: []TaxBracket{
			{Threshold: 200000, Rate: 35}, // 35% for income above $200K
			{Threshold: 100000, Rate: 25}, // 25% for income above $100K
			{Threshold: 50000, Rate: 15},  // 15% for income above $50K
			{Threshold: 20000, Rate: 5},   // 5% for income above $20K
		},

		ReserveValues:      []int{reserves},
		CollectedTaxValues: []int{0},
	}
}
