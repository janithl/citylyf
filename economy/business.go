package economy

import (
	"citylyf/entities"
	"math"
	"math/rand"
)

// Company represents a business entity with jobs
type Company struct {
	Name        string
	Industry    string
	Revenue     float64
	Expenses    float64
	LastProfit  float64
	JobOpenings map[entities.CareerLevel]int // Available job positions at each level
}

// CalculateProfit computes net profit after taxes
func (c *Company) CalculateProfit(m Market) float64 {
	taxedAmount := c.Revenue * m.CorporateTax
	c.LastProfit = c.Revenue - (c.Expenses + taxedAmount)
	return c.LastProfit
}

// DetermineJobOpenings calculates jobs available based on economic factors
func (c *Company) DetermineJobOpenings(m Market) {
	baseJobs := map[entities.CareerLevel]int{
		entities.EntryLevel:     rand.Intn(10) + 5, // 5-15 jobs
		entities.MidLevel:       rand.Intn(5) + 2,  // 2-7 jobs
		entities.SeniorLevel:    rand.Intn(3) + 1,  // 1-4 jobs
		entities.ExecutiveLevel: rand.Intn(2),      // 0-1 jobs
	}

	// Adjust based on economic conditions
	marketMultiplier := 1.0

	// Interest Rate Effect: High rates slow down hiring
	if m.InterestRate > 5 {
		marketMultiplier -= 0.3
	}

	// Inflation Effect: High inflation discourages hiring
	if m.LastInflationRate > 6 {
		marketMultiplier -= 0.2
	}

	// Government Spending Effect: More spending stimulates job creation
	if m.GovernmentSpending > 5 {
		marketMultiplier += 0.2
	}

	// Market Sentiment Effect: High confidence = More job openings
	marketMultiplier += m.LastMarketSentiment * 0.1

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
