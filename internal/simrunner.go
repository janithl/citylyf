package internal

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/janithl/citylyf/internal/economy"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/gamefile"
	"github.com/janithl/citylyf/internal/people"
)

type SimRunner struct {
	employment         *economy.Employment
	calculationService *economy.CalculationService
	ticker             *time.Ticker
	done               chan bool
}

func (sr *SimRunner) NewGame(gamePath *string) {
	if gamePath != nil && gamefile.CheckExists(*gamePath) { // load sim from savegame file
		gamefile.Load(*gamePath)
	} else { // create a new simulation
		entities.StartNewSim()
	}

	sr.employment = &economy.Employment{CompanyService: &economy.CompanyService{}}
	sr.calculationService = economy.NewCalculationService(sr.employment.CompanyService)

	if len(entities.Sim.Companies) == 0 {
		// set up some initial entities.Sim.Companies
		for i := 0; i < 8+rand.Intn(8); i++ {
			entities.Sim.Mutex.Lock()
			newCompany := sr.employment.CompanyService.GenerateRandomCompany(entities.GetRandomCompanySize(), entities.GetRandomIndustry())
			entities.Sim.Companies.Add(newCompany)
			entities.Sim.Mutex.Unlock()
			fmt.Printf("[ Econ ] %s (%s) founded!\n", newCompany.Name, newCompany.Industry)
		}
	}

	sr.ticker = time.NewTicker(100 * time.Millisecond) // tick every 1/10th of a second
	sr.done = make(chan bool)                          // channel to send kill signal to goroutine
}

func (sr *SimRunner) GameTick() {
	entities.Sim.Houses.PlaceHousing()
	people.Immigrate()
	sr.employment.AssignJobs()
	people.Emigrate()
	people.SimulateLifecycle()
	entities.Sim.Market.ReviseInterestRate()
	sr.calculationService.CalculateEconomy()
}

func (sr *SimRunner) RunGameLoop() {
	for {
		select {
		case <-sr.done:
			return
		case <-sr.ticker.C:
			entities.Sim.Mutex.Lock()
			if entities.Sim.SimulationSpeed != entities.Pause {
				entities.Sim.Tick(sr.GameTick)
				entities.Sim.SendStats()
			}
			entities.Sim.Mutex.Unlock()
		}
	}
}

func (sr *SimRunner) EndGame() {
	if sr.ticker != nil {
		sr.ticker.Stop()
	}
	if sr.done != nil {
		sr.done <- true
	}
}

func (sr *SimRunner) SaveGame(gamePath *string) {
	if gamePath != nil {
		gamefile.Save(*gamePath)
	}
}
