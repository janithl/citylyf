package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/janithl/citylyf/internal/economy"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/people"
	"github.com/janithl/citylyf/internal/ui"
)

func main() {
	savePathPtr := flag.String("savePath", "", "where should we load/save the game?")
	flag.Parse()

	if *savePathPtr != "" && checkFileExists(*savePathPtr) { // load sim from savegame file
		loadGame(*savePathPtr)
	} else { // create a new simulation
		entities.Sim = entities.NewSimulation(2020, 1000000)
	}

	employment := economy.Employment{CompanyService: &economy.CompanyService{}}
	calculationService := economy.NewCalculationService(employment.CompanyService)

	if len(entities.Sim.Companies) == 0 {
		// set up some initial entities.Sim.Companies
		for i := 0; i < 8+rand.Intn(8); i++ {
			entities.Sim.Mutex.Lock()
			newCompany := employment.CompanyService.GenerateRandomCompany(entities.GetRandomCompanySize(), entities.GetRandomIndustry())
			entities.Sim.Companies.Add(newCompany)
			entities.Sim.Mutex.Unlock()
			fmt.Printf("[ Econ ] %s (%s) founded!\n", newCompany.Name, newCompany.Industry)
		}
	}

	ticker := time.NewTicker(100 * time.Millisecond) // tick every 1/10th of a second
	done := make(chan bool)                          // channel to send kill signal to goroutine

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				entities.Sim.Mutex.Lock()
				if entities.Sim.SimulationSpeed != entities.Pause {
					entities.Sim.Tick(func() {
						entities.Sim.Houses.PlaceHousing()
						people.Immigrate()
						employment.AssignJobs()
						people.Emigrate()
						entities.Sim.Market.ReviseInterestRate()
						calculationService.CalculateEconomy()
					})
					entities.Sim.SendStats()
				}
				entities.Sim.Mutex.Unlock()
			}
		}
	}()

	// run simulation until UI is closed
	ui.RunGame()
	ticker.Stop()
	done <- true

	if *savePathPtr != "" {
		saveGame(*savePathPtr)
	}
}
