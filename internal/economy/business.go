package economy

import (
	"fmt"
	"math/rand"

	"github.com/janithl/citylyf/internal/entities"
)

type CompanyService struct {
	LastCompanyID int
}

func NewCompanyService() *CompanyService {
	return &CompanyService{
		LastCompanyID: 1000, // start IDs from 1000
	}
}

// GenerateRandomCompany creates a company with random industry and financials
func (c *CompanyService) GenerateRandomCompany() *entities.Company {
	// Assign financials based on industry type
	baseRevenue := rand.Float64()*5_000_000 + 1_000_000 // Revenue between $1M - $6M
	expenseRatio := rand.Float64()*0.4 + 0.5            // Expenses are 50-90% of revenue
	expenses := baseRevenue * expenseRatio

	// Generate a random company name
	companyNames := []string{"Global", "NextGen", "Quantum", "Future", "Vertex", "Synergy", "Omni", "Pinnacle", "Apex", "Horizon"}
	companySuffix := []string{"Corp", "Industries", "Systems", "Group", "Technologies", "Enterprises"}
	companyName := fmt.Sprintf("%s %s", companyNames[rand.Intn(len(companyNames))], companySuffix[rand.Intn(len(companySuffix))])

	company := entities.Company{
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
	return &company
}

// AddCompany adds a new company
func (c *CompanyService) AddCompany(company *entities.Company) {
	c.LastCompanyID += 1
	company.ID = c.LastCompanyID
	entities.Sim.Companies[company.ID] = company
}

// RemoveCompany removes a company
func (c *CompanyService) RemoveCompany(companyID int) {
	delete(entities.Sim.Companies, companyID)
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
