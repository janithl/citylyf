package economy

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"time"

	"github.com/janithl/citylyf/internal/entities"
)

type CalculationService struct {
	companyService  *CompanyService
	lastCalculation time.Time
	nextCalculation time.Time
}

func NewCalculationService(cs *CompanyService) *CalculationService {
	return &CalculationService{
		companyService:  cs,
		lastCalculation: entities.Sim.Date,
		nextCalculation: entities.Sim.Date.AddDate(0, 1, 0),
	}
}

func (cs *CalculationService) CalculateEconomy() {
	if entities.Sim.Date.Before(cs.nextCalculation) { // run monthly
		return
	}
	daysSinceLastCalculation := entities.Sim.Date.Sub(cs.lastCalculation).Hours() / entities.HoursPerDay
	cs.lastCalculation = entities.Sim.Date
	cs.nextCalculation = cs.lastCalculation.AddDate(0, 1, 0)

	// calculate impact of population growth on city economy
	populationGrowth := entities.Sim.People.PopulationGrowthRate()

	entities.Sim.People.UpdatePopulationValues()
	entities.Sim.People.CalculateAgeGroups()
	entities.Sim.People.CalculateUnemployment()

	entities.Sim.Market.CalculateInflation(populationGrowth)
	marketGrowth := entities.Sim.Market.CalculateMarketGrowth()
	entities.Sim.Market.CalculateHousingAndRetailDemand(len(entities.Sim.Houses), entities.Sim.Houses.GetFreeHouses())
	entities.Sim.Market.UpdateMarketValue(marketGrowth)

	fmt.Printf("[ Econ ] %s | Next calculation on %s\n", entities.Sim.GetStats(), cs.nextCalculation.Format("2006-01-02"))

	if marketGrowth > 0 && rand.IntN(100) < 5 { // 5% chance of a farm being opened during good times
		newFarm := cs.companyService.GenerateRandomCompany(entities.SME, entities.Agriculture)
		entities.Sim.Companies.PlaceAgriculture(newFarm)
		fmt.Printf("[ Econ ] Growth! %s (%s) founded!\n", newFarm.Name, newFarm.Industry)
	} else if entities.Sim.Market.RetailDemand > 0.01 && rand.IntN(100) < 25 { // 25% chance of a shop being opened when retail demand over 1%
		newRetailCompany := cs.companyService.GenerateRandomCompany(entities.Micro, entities.Retail)
		entities.Sim.Companies.PlaceRetail(newRetailCompany)
		fmt.Printf("[ Econ ] Growth! %s (%s) founded!\n", newRetailCompany.Name, newRetailCompany.Industry)
	}

	totalProfits := 0.0
	for id, company := range entities.Sim.Companies {
		totalProfits += company.CalculateProfit(daysSinceLastCalculation)
		company.DetermineJobOpenings()
		entities.Sim.Companies[id] = company
	}
	entities.Sim.Market.ReportCompanyProfits(totalProfits)

	// do govt interest calcuations (monthly)
	monthlyInterestRate := (entities.Sim.Market.InterestRate() / 100) * (daysSinceLastCalculation / entities.DaysPerYear)
	entities.Sim.Government.Reserves += int(float64(entities.Sim.Government.Reserves) * monthlyInterestRate)

	// calculate monthly pay and interest for households
	for household := range maps.Values(entities.Sim.People.Households) {
		household.CalculateMonthlyBudget(cs.companyService.AddPayToPayroll)
		household.Savings += int(float64(household.Savings) * monthlyInterestRate)
	}

	// collect taxes, revise rents and calculate regional stats and sales
	entities.Sim.Government.CollectTaxes()
	entities.Sim.People.UpdateAverageWageValues()
	entities.Sim.Houses.ReviseRents()
	entities.Sim.Geography.Regions.CalculateRegionalStats()
	entities.Sim.Geography.Regions.CalculateRegionalSales()
}
