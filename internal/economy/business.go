package economy

import (
	"fmt"
	"math/rand"

	"github.com/janithl/citylyf/internal/entities"
)

// GenerateRandomCompany creates a company with random industry and financials
func GenerateRandomCompany() entities.Company {
	// Assign financials based on industry type
	baseRevenue := rand.Float64()*5_000_000 + 1_000_000 // Revenue between $1M - $6M
	expenseRatio := rand.Float64()*0.4 + 0.5            // Expenses are 50-90% of revenue
	expenses := baseRevenue * expenseRatio

	// Generate a random company name
	companyNames := []string{"Global", "NextGen", "Quantum", "Future", "Vertex", "Synergy", "Omni", "Pinnacle", "Apex", "Horizon"}
	companySuffix := []string{"Corp", "Industries", "Systems", "Group", "Technologies", "Enterprises"}
	companyName := fmt.Sprintf("%s %s", companyNames[rand.Intn(len(companyNames))], companySuffix[rand.Intn(len(companySuffix))])

	company := entities.Company{
		ID:           rand.Intn(999) + 1000,
		Name:         companyName,
		Industry:     entities.GetRandomIndustry(),
		FoundingDate: entities.Sim.Date,
		JobOpenings:  make(map[entities.CareerLevel]int),
		LastRevenue:  baseRevenue,
		LastExpenses: expenses,
		LastProfit:   baseRevenue - expenses,
	}

	company.CalculateProfit()
	return company
}
