package entities

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

type Companies map[int]*Company

// Add adds a new company
func (c Companies) Add(company *Company) {
	company.ID = Sim.GetNextID()
	c[company.ID] = company
}

// Remove removes a company
func (c Companies) Remove(companyID int) {
	delete(c, companyID)
}

// GetIDs returns a sorted list of company IDs
func (c Companies) GetIDs() []int {
	IDs := []int{}
	for company := range maps.Values(c) {
		IDs = append(IDs, company.ID)
	}
	slices.Sort(IDs)
	return IDs
}

func (c Companies) PlaceRetail(newCompany *Company) {
	for x := 0; x < Sim.Geography.Size; x++ {
		for y := 0; y < Sim.Geography.Size; y++ {
			if !Sim.Geography.tiles[x][y].Shop && Sim.Geography.tiles[x][y].Zone == RetailZone {
				Sim.Geography.tiles[x][y].Zone = NoZone

				roadDir := Sim.Geography.getAccessRoad(x, y)
				if roadDir == "" {
					return
				}

				Sim.Geography.tiles[x][y].Shop = true
				newCompany.Location = Point{X: x, Y: y}
				newCompany.RoadDirection = roadDir
				c.Add(newCompany)

				return
			}
		}
	}
}

func (c Companies) GetLocationCompany(x, y int) *Company {
	for company := range maps.Values(c) {
		if company.Location.X == x && company.Location.Y == y {
			return company
		}
	}
	return nil
}

// Company represents a business entity with jobs
type Company struct {
	ID            int
	Name          string
	Industry      Industry
	CompanySize   CompanySize
	Location      Point
	RoadDirection Direction
	FoundingDate  time.Time
	JobOpenings   map[CareerLevel]int // Available job positions at each level
	Employees     []int               // Employee IDs
	RetailSales   float64
	TaxPayable    float64
	FixedCosts    float64
	Payroll       float64

	// Historical
	LastRevenue, LastExpenses, LastProfit float64
}

// CalculateProfit computes monthly net profit
func (c *Company) CalculateProfit(monthLength float64) float64 {
	lastMarketGrowthRate := utils.GetLastValue(Sim.Market.History.MarketGrowthRate)

	// **Monthly Expense Growth**: Inflation applied proportionally
	inflationMultiplier := 1.0 + (Sim.Market.InflationRate() / 2400) // Divided by 2400 for smoother monthly change

	// **Cost-cutting for struggling companies**: Reduces expenses if past profits were negative
	if c.LastProfit < 0 {
		inflationMultiplier *= 0.95 // Expenses grow slower for struggling businesses
	}
	c.FixedCosts *= inflationMultiplier // Adjust fixed costs with inflation
	c.LastExpenses = c.FixedCosts + c.Payroll
	c.Payroll = 0.0 // Reset payroll liabilites

	if c.Industry == Retail { // For retail, revenue == sales TODO: Fix negative
		c.LastRevenue = c.RetailSales
	} else {
		// **Monthly Revenue Growth**: Adjusts based on market conditions for non-retail
		revenueMultiplier := 1.0 + (lastMarketGrowthRate / 1200) // Gradual revenue increase
		if c.LastProfit > 0 {
			revenueMultiplier += 0.002 // Small bonus growth for profitable companies
		}
		c.LastRevenue *= revenueMultiplier
	}

	// **Calculate Profit**
	grossProfit := c.LastRevenue - c.LastExpenses

	// **Apply Corporate Tax (Adjusted for Monthly Periods)**
	if grossProfit > 0 {
		monthlyTaxRate := Sim.Government.CorporateTaxRate * monthLength / DaysPerYear // Convert annual tax rate to monthly profit cycle
		taxedAmount := math.Ceil(grossProfit * (monthlyTaxRate / 100.0))              // round to nearest dollar
		c.LastProfit = grossProfit - taxedAmount

		// Store unpaid tax in liability account
		c.TaxPayable += taxedAmount
	} else {
		c.LastProfit = grossProfit
	}

	// **Loss Limiter**: Prevents companies from collapsing too fast
	if c.LastProfit < -c.LastRevenue*0.25 { // Reduced from 50% to 25% for monthly scaling
		c.LastProfit = -c.LastRevenue * 0.25
	}

	return c.LastProfit
}

// GetNumberOfJobOpenings returns the number of job openings
func (c *Company) GetNumberOfJobOpenings() int {
	openings := 0
	for i := 0; i < len(CareerLevels); i++ {
		openings += c.JobOpenings[CareerLevels[i]]
	}
	return openings
}

// GetNumberOfEmployees returns the total number of employees
func (c *Company) GetNumberOfEmployees() int {
	return len(c.Employees)
}

// GetEmployees returns a list of employees
func (c *Company) GetEmployees() []*Person {
	employees := []*Person{}
	for _, employeeID := range c.Employees {
		employees = append(employees, Sim.People.People[employeeID])
	}
	return employees
}

// RemoveEmployee removes an employee from the company
func (c *Company) RemoveEmployee(employeeID int) {
	c.Employees = slices.DeleteFunc(c.Employees, func(id int) bool {
		return id == employeeID
	})
}

// DetermineJobOpenings calculates jobs available based on economic factors
func (c *Company) DetermineJobOpenings() {
	lastMarketSentiment := utils.GetLastValue(Sim.Market.History.MarketSentiment)
	baseJobs := c.CompanySize.GetBaseJobs()

	// Adjust based on economic conditions
	marketMultiplier := 1.0

	// Interest Rate Effect: High rates slow down hiring
	if Sim.Market.InterestRate() > 5 {
		marketMultiplier -= 0.3
	}

	// Inflation Effect: High inflation discourages hiring
	if Sim.Market.InflationRate() > 6 {
		marketMultiplier -= 0.2
	}

	// Government Spending Effect: More spending stimulates job creation
	if Sim.Government.GetGovernmentSpending() > 5 {
		marketMultiplier += 0.2
	}

	// Market Sentiment Effect: High confidence = More job openings
	marketMultiplier += lastMarketSentiment * 0.1

	// Adjust hiring based on profitability
	if c.LastProfit < 0 {
		marketMultiplier -= 0.5 // If losing money, reduce hiring
	}

	// Apply adjustments
	for level, jobs := range baseJobs {
		adjustedJobs := int(math.Round(float64(jobs) * marketMultiplier))
		if openings, exists := c.JobOpenings[level]; exists && openings > 0 { // if job openings already exists, use that value
			adjustedJobs = int(math.Round(float64(openings) * marketMultiplier))
		}
		if adjustedJobs < 1 {
			adjustedJobs = 1 // Prevent zero or negative jobs
		}
		c.JobOpenings[level] = adjustedJobs
	}
}

func (c *Company) CompanyAge() int {
	duration := Sim.Date.Sub(c.FoundingDate)
	return int(duration.Hours() / HoursPerYear)
}

func (c *Company) GetID() int {
	return c.ID
}

func (c *Company) GetStats() string {
	return fmt.Sprintf("%5d %-25s %-5s %02d/%02d %4d %-18s %-10s", c.ID, c.Name, c.CompanySize, c.GetNumberOfEmployees(), c.GetNumberOfJobOpenings(), c.FoundingDate.Year(), c.Industry, utils.FormatCurrency(c.LastProfit, "$"))
}
