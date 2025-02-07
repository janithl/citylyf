package entities

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

// Company represents a business entity with jobs
type Company struct {
	ID           int
	Name         string
	Industry     Industry
	FoundingDate time.Time
	JobOpenings  map[CareerLevel]int // Available job positions at each level

	// Historical
	LastRevenue  float64
	LastExpenses float64
	LastProfit   float64
}

// CalculateProfit computes net profit every 28 days
func (c *Company) CalculateProfit() float64 {
	lastInflationRate := utils.GetLastValue(Sim.Market.History.InflationRate)
	lastMarketGrowthRate := utils.GetLastValue(Sim.Market.History.MarketGrowthRate)

	// **Monthly Expense Growth**: Inflation applied proportionally
	inflationMultiplier := 1.0 + (lastInflationRate / 2400) // Divided by 2400 for smoother monthly change

	// **Cost-cutting for struggling companies**: Reduces expenses if past profits were negative
	if c.LastProfit < 0 {
		inflationMultiplier *= 0.95 // Expenses grow slower for struggling businesses
	}
	c.LastExpenses *= inflationMultiplier // Adjust expenses

	// **Monthly Revenue Growth**: Adjusts based on market conditions
	revenueMultiplier := 1.0 + (lastMarketGrowthRate / 1200) // Gradual revenue increase
	if c.LastProfit > 0 {
		revenueMultiplier += 0.002 // Small bonus growth for profitable companies
	}
	c.LastRevenue *= revenueMultiplier

	// **Calculate Profit**
	grossProfit := c.LastRevenue - c.LastExpenses
	if grossProfit > 0 {
		taxedAmount := grossProfit * (Sim.Market.CorporateTax / 100.0)
		c.LastProfit = grossProfit - taxedAmount
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

// DetermineJobOpenings calculates jobs available based on economic factors
func (c *Company) DetermineJobOpenings() {
	lastInflationRate := utils.GetLastValue(Sim.Market.History.InflationRate)
	lastMarketSentiment := utils.GetLastValue(Sim.Market.History.MarketSentiment)

	baseJobs := map[CareerLevel]int{
		EntryLevel:     rand.Intn(10) + 5, // 5-15 jobs
		MidLevel:       rand.Intn(5) + 2,  // 2-7 jobs
		SeniorLevel:    rand.Intn(3) + 1,  // 1-4 jobs
		ExecutiveLevel: rand.Intn(2),      // 0-1 jobs
	}

	// Adjust based on economic conditions
	marketMultiplier := 1.0

	// Interest Rate Effect: High rates slow down hiring
	if Sim.Market.InterestRate > 5 {
		marketMultiplier -= 0.3
	}

	// Inflation Effect: High inflation discourages hiring
	if lastInflationRate > 6 {
		marketMultiplier -= 0.2
	}

	// Government Spending Effect: More spending stimulates job creation
	if Sim.Market.GovernmentSpending > 5 {
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
		if adjustedJobs < 0 {
			adjustedJobs = 0 // Prevent negative jobs
		}
		c.JobOpenings[level] = adjustedJobs
	}
}

func (c *Company) CompanyAge() int {
	duration := Sim.Date.Sub(c.FoundingDate)
	return int(duration.Hours() / HoursPerYear)
}

func (c *Company) GetStats() string {
	return fmt.Sprintf("%4d %-28s %d %-18s %-10s", c.ID, c.Name, c.FoundingDate.Year(), c.Industry, utils.FormatCurrency(c.LastProfit, "$"))
}
