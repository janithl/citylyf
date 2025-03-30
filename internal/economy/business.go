package economy

import (
	"math/rand/v2"

	"github.com/janithl/citylyf/internal/entities"
)

type CompanyService struct{}

// GenerateRandomCompany creates a company with random industry and financials
func (c *CompanyService) GenerateRandomCompany(companySize entities.CompanySize, industry entities.Industry) *entities.Company {
	// Assign financials based on industry type
	baseRevenue := companySize.GetBaseRevenue()
	expenseRatio := rand.Float64()*0.4 + 0.5 // Expenses are 50-90% of revenue
	expenses := baseRevenue * expenseRatio

	company := entities.Company{
		Name:             entities.Sim.NameService.GetCompanyName(),
		Industry:         industry,
		CompanySize:      companySize,
		FoundingDate:     entities.Sim.Date,
		NextWageRevision: entities.Sim.Date.AddDate(1, 0, 0),
		JobOpenings:      make(map[entities.CareerLevel]int),
		LastRevenue:      baseRevenue,
		LastExpenses:     expenses,
		FixedCosts:       expenses,
		Payroll:          0,
		LastProfit:       baseRevenue - expenses,
	}

	company.CalculateProfit(31)
	company.DetermineJobOpenings()
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

// AddPayToPayroll adds your payroll payment as a liability to the company
func (c *CompanyService) AddPayToPayroll(companyID int, payAmount float64) {
	company, ok := entities.Sim.Companies[companyID]
	if ok {
		company.Payroll -= payAmount
		entities.Sim.Companies[companyID] = company
	}
}
