package entities_test

import (
	"testing"

	"github.com/janithl/citylyf/internal/economy"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/people"
)

func TestReviseWages(t *testing.T) {
	entities.Sim = entities.NewSimulation(2020, 1e6)
	entities.Sim.SimulationSpeed = entities.Fast

	employment := economy.Employment{CompanyService: &economy.CompanyService{}}
	calculationService := economy.NewCalculationService(employment.CompanyService)
	industry := entities.GetRandomIndustry()

	newCompany := employment.CompanyService.GenerateRandomCompany(entities.Large, industry)
	entities.Sim.Companies.Add(newCompany)

	for range 10 {
		household := people.CreateHousehold()
		entities.Sim.People.Households[household.ID] = household
		for _, person := range household.GetMembers() {
			if person.IsEmployable() {
				newCompany.AddEmployee(person.ID)
				person.EmployerID = newCompany.ID
			}
		}
	}

	initialAvgWage := entities.Sim.People.AverageWage()
	for range 366 * 5 {
		entities.Sim.Tick(func() {
			employment.AssignJobs()
			people.SimulateLifecycle()
			entities.Sim.Market.ReviseInterestRate()
			calculationService.CalculateEconomy()
		})
		if initialAvgWage == 0 {
			initialAvgWage = entities.Sim.People.AverageWage()
		}
	}

	finalAvgWage := entities.Sim.People.AverageWage()
	expectedMaxWageGrowth := 1.28  // %5 compounding over 5 years
	expectedMinWageGrowth := 1.025 // %0.5 compounding over 5 years
	actualWageGrowth := finalAvgWage / initialAvgWage

	if actualWageGrowth > expectedMaxWageGrowth {
		t.Errorf(`ReviseWages(): Expected less than %.2f growth over 5 years at max 5%% annually, got %.2f\n`, expectedMaxWageGrowth, actualWageGrowth)
	} else if actualWageGrowth < expectedMinWageGrowth {
		t.Errorf(`ReviseWages(): Expected more than %.2f growth over 5 years at min 0.5%% annually, got %.2f\n`, expectedMinWageGrowth, actualWageGrowth)
	}
}
