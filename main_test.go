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
	entities.Sim = entities.NewSimulation(2020, 100000)
	employment := economy.Employment{CompanyService: &economy.CompanyService{}}
	peopleService := people.NewPeopleService()
	calculationService := economy.NewCalculationService(employment.CompanyService)
	entities.Sim.SimulationSpeed = entities.Slow

	for i := 0; i < 10; i++ {
		newCompany := employment.CompanyService.GenerateRandomCompany()
		entities.Sim.Companies.Add(newCompany)
	}

	simSize := entities.Sim.Geography.Size
	for i := 0; i < 50; i++ {
		x, y := rand.IntN(simSize), rand.IntN(simSize)
		entities.PlaceRoad(x-1, y, x+1, y, entities.Asphalt)
		entities.Sim.Houses.AddHouse(x, y+1, 2+rand.IntN(3))
	}

	for i := 0; i < b.N; i++ {
		entities.Sim.Tick()
		entities.Sim.People.MoveIn(peopleService.CreateHousehold)
		employment.AssignJobs()
		entities.Sim.People.MoveOut(employment.CompanyService.RemoveEmployeeFromCompany)
		entities.Sim.Market.ReviseInterestRate()
		calculationService.CalculateEconomy()
		entities.Sim.SendStats()
	}

	for _, household := range entities.Sim.People.Households {
		fmt.Println(household.GetStats())
		fmt.Println(household.GetMemberStats())
	}
}
