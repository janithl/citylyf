package entities

import (
	"fmt"
	"maps"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

type Government struct {
	Reserves, CapEx                int
	LastCalculationYear            int
	CorporateTaxRate, SalesTaxRate float64              // Flat corporate tax rate and sales tax
	IncomeTaxBrackets              []TaxBracket         // Progressive income tax brackets
	Expenses                       map[CostType]float64 // Holds goverment expenses

	// Historical values
	ReserveValues, IncomeValues, CapExValues, OpExValues []int
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

	// Collect sales and corporate tax and reset tax payable account
	salesTaxesCollected := 0
	corporateTaxesCollected := 0
	for id := range Sim.Companies {
		corporateTaxesCollected += int(Sim.Companies[id].CorpTaxPayable)
		salesTaxesCollected += int(Sim.Companies[id].SalesTaxPayable)
		Sim.Companies[id].CorpTaxPayable = 0.0
		Sim.Companies[id].SalesTaxPayable = 0.0
	}
	fmt.Printf("[  Tax ] Collected $%d in sales taxes, and $%d corporate taxes\n", salesTaxesCollected, corporateTaxesCollected)

	// add collected taxes to government income
	totalTaxesCollected := personalTaxesCollected + corporateTaxesCollected + salesTaxesCollected
	g.IncomeValues = utils.AddFifo(g.IncomeValues, totalTaxesCollected, 10)

	// get opex
	opEx := g.CalculateOpEx()
	g.OpExValues = utils.AddFifo(g.OpExValues, opEx, 10)

	// calculate final reserves
	g.Reserves = g.Reserves + totalTaxesCollected - g.CapEx - opEx
	g.ReserveValues = utils.AddFifo(g.ReserveValues, g.Reserves, 10)
	fmt.Printf("[  Tax ] %d: Income: $%d, CapEx: $%d, OpEx: $%d, Total Government Reserves: $%d\n",
		g.LastCalculationYear, totalTaxesCollected, g.CapEx, opEx, g.Reserves)

	// revise government expenses
	g.ReviseExpenses()

	// reset capex spend
	g.LastCalculationYear = Sim.Date.Year()
	g.CapExValues = utils.AddFifo(g.CapExValues, g.CapEx, 10)
	g.CapEx = 0
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

func (g *Government) GetReservesAtHand() float64 {
	return float64(g.Reserves - g.CapEx)
}

// NewGovernment initializes the government system with reserves and progressive tax brackets
func NewGovernment(reserves int, startDate time.Time) *Government {
	return &Government{
		Reserves:            reserves,
		CapEx:               0,
		LastCalculationYear: startDate.Year(),
		CorporateTaxRate:    9.5,
		SalesTaxRate:        12.5,
		IncomeTaxBrackets: []TaxBracket{
			{Threshold: 200000, Rate: 35}, // 35% for income above $200K
			{Threshold: 100000, Rate: 25}, // 25% for income above $100K
			{Threshold: 50000, Rate: 15},  // 15% for income above $50K
			{Threshold: 20000, Rate: 5},   // 5% for income above $20K
		},
		Expenses:      NewExpenses(),
		ReserveValues: []int{reserves},
		IncomeValues:  []int{0},
		CapExValues:   []int{0},
		OpExValues:    []int{0},
	}
}
