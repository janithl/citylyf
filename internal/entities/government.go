package entities

import "fmt"

type Government struct {
	Reserves            int
	ReserveValues       []int // Historical values
	LastCalculationYear int
	CorporateTaxRate    float64      // Flat corporate tax rate
	IncomeTaxBrackets   []TaxBracket // Progressive income tax brackets
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

	totalCollected := 0

	// **Collect Household Income Taxes**
	for i := range Sim.People.Households {
		household := &Sim.People.Households[i]
		householdTax := g.CalculateIncomeTax(household.AnnualIncome())

		// Deduct tax from household wealth
		household.Savings -= householdTax
		g.Reserves += householdTax
		totalCollected += householdTax
	}

	fmt.Printf("[  Tax ] Collected $%d in taxes for %d. Government reserves: $%d\n",
		totalCollected, g.LastCalculationYear, g.Reserves)
	g.LastCalculationYear = Sim.Date.Year()
	g.updateReserveValues()
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

// Append current reserve value to history
func (g *Government) updateReserveValues() {
	if len(g.ReserveValues) >= 10 {
		g.ReserveValues = g.ReserveValues[1:] // Remove first element (FIFO behavior)
	}
	g.ReserveValues = append(g.ReserveValues, g.Reserves)
}

// NewGovernment initializes the government system with reserves and progressive tax brackets
func NewGovernment(reserves int) *Government {
	return &Government{
		Reserves:            reserves,
		ReserveValues:       []int{reserves},
		LastCalculationYear: Sim.Date.Year(),
		CorporateTaxRate:    5.0,
		IncomeTaxBrackets: []TaxBracket{
			{Threshold: 200000, Rate: 35}, // 35% for income above $200K
			{Threshold: 100000, Rate: 25}, // 25% for income above $100K
			{Threshold: 50000, Rate: 15},  // 15% for income above $50K
			{Threshold: 20000, Rate: 5},   // 5% for income above $20K
		},
	}
}
