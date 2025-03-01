package main

import (
	"encoding/json"
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
	jsonPtr := flag.Bool("json", false, "should output be in json?")
	flag.Parse()

	entities.Sim = entities.NewSimulation(2020, 100000)
	employment := economy.Employment{CompanyService: economy.NewCompanyService()}
	peopleService := people.NewPeopleService()
	calculationService := economy.NewCalculationService(employment.CompanyService)

	// set up some initial entities.Sim.Companies
	for i := 0; i < 4+rand.Intn(4); i++ {
		newCompany := employment.CompanyService.GenerateRandomCompany()
		employment.CompanyService.AddCompany(newCompany)
		fmt.Printf("[ Econ ] %s (%s) founded!\n", newCompany.Name, newCompany.Industry)
	}

	ticker := time.NewTicker(100 * time.Millisecond) // tick every 1/10th of a second
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if entities.Sim.SimulationSpeed != entities.Pause {
					entities.Sim.Mutex.Lock()
					entities.Sim.Tick()
					entities.Sim.People.MoveIn(peopleService.CreateHousehold)
					employment.AssignJobs()
					entities.Sim.People.MoveOut(employment.CompanyService.RemoveEmployeeFromCompany)
					entities.Sim.Market.ReviseInterestRate()
					calculationService.CalculateEconomy()
					entities.Sim.Mutex.Unlock()
				}
			}
		}
	}()

	// run simulation until UI is closed
	ui.RunGame()
	ticker.Stop()
	done <- true

	if *jsonPtr {
		printFinalState()
	}
}

func printFinalState() {
	peopleJson, err := json.Marshal(entities.Sim.People)
	if err != nil {
		fmt.Println(err)
		return
	}
	compJson, err := json.Marshal(entities.Sim.Companies)
	if err != nil {
		fmt.Println(err)
		return
	}
	houseJson, err := json.Marshal(entities.Sim.Houses)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("[ JSON ] Population: ", string(peopleJson))
	fmt.Println("[ JSON ] Companies: ", string(compJson))
	fmt.Println("[ JSON ] Houses: ", string(houseJson))
}
