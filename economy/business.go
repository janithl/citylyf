package economy

import (
	"citylyf/entities"
	"fmt"
	"math"
	"math/rand"
)

// Company represents a business entity with jobs
type Company struct {
	ID          int
	Name        string
	Industry    entities.Industry
	JobOpenings map[entities.CareerLevel]int // Available job positions at each level

	// Historical
	LastRevenue  float64
	LastExpenses float64
	LastProfit   float64
}

// CalculateProfit computes net profit after taxes
func (c *Company) CalculateProfit(m Market) float64 {
	c.LastExpenses += c.LastExpenses * (m.LastInflationRate / 100)
	c.LastRevenue += c.LastRevenue * (m.MarketGrowth() / 100)

	taxedAmount := c.LastRevenue * (m.CorporateTax / 100.0)
	c.LastProfit = c.LastRevenue - (c.LastExpenses + taxedAmount)
	return c.LastProfit
}

// GetNumberOfJobOpenings returns the number of job openings
func (c *Company) GetNumberOfJobOpenings() int {
	openings := 0
	for i := 0; i < len(entities.CareerLevels); i++ {
		openings += c.JobOpenings[entities.CareerLevels[i]]
	}
	return openings
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

// GenerateRandomCompany creates a company with random industry and financials
func GenerateRandomCompany(m Market) Company {
	// Assign financials based on industry type
	baseRevenue := rand.Float64()*5_000_000 + 1_000_000 // Revenue between $1M - $6M
	expenseRatio := rand.Float64()*0.4 + 0.5            // Expenses are 50-90% of revenue
	expenses := baseRevenue * expenseRatio

	// Generate a random company name
	companyNames := []string{"Global", "NextGen", "Quantum", "Future", "Vertex", "Synergy", "Omni", "Pinnacle", "Apex", "Horizon"}
	companySuffix := []string{"Corp", "Industries", "Systems", "Group", "Technologies", "Enterprises"}
	companyName := fmt.Sprintf("%s %s", companyNames[rand.Intn(len(companyNames))], companySuffix[rand.Intn(len(companySuffix))])

	company := Company{
		ID:           rand.Intn(999) + 1000,
		Name:         companyName,
		Industry:     entities.GetRandomIndustry(),
		JobOpenings:  make(map[entities.CareerLevel]int),
		LastRevenue:  baseRevenue,
		LastExpenses: expenses,
		LastProfit:   baseRevenue - expenses,
	}
	company.CalculateProfit(m)

	return company
}
