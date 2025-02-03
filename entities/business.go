package entities

import (
	"math"
	"math/rand"
)

// Company represents a business entity with jobs
type Company struct {
	ID          int
	Name        string
	Industry    Industry
	JobOpenings map[CareerLevel]int // Available job positions at each level

	// Historical
	LastRevenue  float64
	LastExpenses float64
	LastProfit   float64
}

// CalculateProfit computes net profit after taxes
func (c *Company) CalculateProfit() float64 {
	c.LastExpenses += c.LastExpenses * (Sim.Market.LastInflationRate / 100)
	c.LastRevenue += c.LastRevenue * (Sim.Market.MarketGrowth() / 100)

	taxedAmount := c.LastRevenue * (Sim.Market.CorporateTax / 100.0)
	c.LastProfit = c.LastRevenue - (c.LastExpenses + taxedAmount)
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
	if Sim.Market.LastInflationRate > 6 {
		marketMultiplier -= 0.2
	}

	// Government Spending Effect: More spending stimulates job creation
	if Sim.Market.GovernmentSpending > 5 {
		marketMultiplier += 0.2
	}

	// Market Sentiment Effect: High confidence = More job openings
	marketMultiplier += Sim.Market.LastMarketSentiment * 0.1

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
