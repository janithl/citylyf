package economy

import (
	"math/rand"

	"github.com/janithl/citylyf/internal/entities"
)

type CompanyService struct{}

// GenerateRandomCompany creates a company with random industry and financials
func (c *CompanyService) GenerateRandomCompany() *entities.Company {
	// Assign financials based on industry type
	baseRevenue := rand.Float64()*5_000_000 + 1_000_000 // Revenue between $1M - $6M
	expenseRatio := rand.Float64()*0.4 + 0.5            // Expenses are 50-90% of revenue
	expenses := baseRevenue * expenseRatio

	company := entities.Company{
		Name:         entities.Sim.NameService.GetCompanyName(),
		Industry:     entities.GetRandomIndustry(),
		FoundingDate: entities.Sim.Date,
		JobOpenings:  make(map[entities.CareerLevel]int),
		LastRevenue:  baseRevenue,
		LastExpenses: expenses,
		FixedCosts:   expenses,
		Payroll:      0,
		LastProfit:   baseRevenue - expenses,
	}

	company.CalculateProfit(31)
	return &company
}

// AddEmployeeToCompany adds your ID to the company list of employees
func (c *CompanyService) AddEmployeeToCompany(companyID int, employeeID int) {
	company, ok := entities.Sim.Companies[companyID]
	if ok {
		company.Employees = append(company.Employees, employeeID)
		entities.Sim.Companies[companyID] = company
	}
}

// RemoveEmployeeFromCompany removes your ID from the company list of employees
func (c *CompanyService) RemoveEmployeeFromCompany(companyID int, employeeID int) {
	company, ok := entities.Sim.Companies[companyID]
	if ok {
		company.RemoveEmployee(employeeID)
		entities.Sim.Companies[companyID] = company
	}
}

// AddPayToPayroll adds your payroll payment as a liability to the company
func (c *CompanyService) AddPayToPayroll(companyID int, payAmount float64) {
	company, ok := entities.Sim.Companies[companyID]
	if ok {
		company.Payroll -= payAmount
		entities.Sim.Companies[companyID] = company
	}
}
