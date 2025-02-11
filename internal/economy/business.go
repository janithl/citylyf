package economy

import (
	"fmt"
	"math/rand"
	"slices"

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
		FixedCosts:   expenses,
		Payroll:      0,
		LastProfit:   baseRevenue - expenses,
	}

	company.CalculateProfit()
	return company
}

// AddEmployeeToCompany adds your ID to the company list of employees
func AddEmployeeToCompany(companyID int, employeeID int) {
	for c := range entities.Sim.Companies {
		if entities.Sim.Companies[c].ID == companyID {
			entities.Sim.Companies[c].Employees = append(entities.Sim.Companies[c].Employees, employeeID)
			return
		}

	}

}

// RemoveEmployeeFromCompany removes your ID from the company list of employees
func RemoveEmployeeFromCompany(companyID int, employeeID int) {
	for c := range entities.Sim.Companies {
		if entities.Sim.Companies[c].ID == companyID {
			entities.Sim.Companies[c].Employees = slices.DeleteFunc(entities.Sim.Companies[c].Employees, func(id int) bool {
				return id == employeeID
			})
			return
		}

	}

}

// AddPayToPayroll adds your payroll payment as a liability to the company
func AddPayToPayroll(companyID int, payAmount float64) {
	for c := range entities.Sim.Companies {
		if entities.Sim.Companies[c].ID == companyID {
			entities.Sim.Companies[c].Payroll -= payAmount
			return
		}

	}

}
