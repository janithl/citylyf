package main_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/janithl/citylyf/internal/economy"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/people"
)

func BenchmarkSim(b *testing.B) {
	entities.Sim = entities.NewSimulation(2020, 1000000)
	employment := economy.Employment{CompanyService: &economy.CompanyService{}}
	calculationService := economy.NewCalculationService(employment.CompanyService)
	entities.Sim.SimulationSpeed = entities.Slow

	for i := 0; i < 10; i++ {
		newCompany := employment.CompanyService.GenerateRandomCompany(entities.GetRandomCompanySize(), entities.GetRandomIndustry())
		entities.Sim.Companies.Add(newCompany)
	}

	simSize := entities.Sim.Geography.Size
	for i := 0; i < 16; i++ {
		x, y := rand.IntN(simSize), rand.IntN(simSize)
		entities.PlaceRoad(entities.Point{X: x - 1, Y: y}, entities.Point{X: x + 1, Y: y}, entities.Asphalt)
		use := entities.ResidentialUse
		if i >= 12 {
			use = entities.RetailUse
		}
		entities.Sim.Geography.PlaceLandUse(entities.Point{X: x - 2, Y: y - 1}, entities.Point{X: x + 2, Y: y + 1}, use)
	}

	for i := 0; i < b.N; i++ {
		entities.Sim.Tick(func() {
			entities.Sim.Houses.PlaceHousing()
			people.Immigrate()
			employment.AssignJobs()
			people.Emigrate()
			people.SimulateLifecycle()
			entities.Sim.Market.ReviseInterestRate()
			calculationService.CalculateEconomy()
		})
		entities.Sim.SendStats()
	}

	for _, household := range entities.Sim.People.Households {
		fmt.Println(household.GetStats())
		fmt.Println(household.GetMemberStats())
	}
}
