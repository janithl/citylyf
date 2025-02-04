package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/janithl/citylyf/economy"
	"github.com/janithl/citylyf/entities"
	"github.com/janithl/citylyf/internal/ui"
	"github.com/janithl/citylyf/people"
)

// TODO
// Household Budgeting - think about rent/mortgage expenses + taxes + savings interest etc
var freeHouses = 100

func main() {
	entities.Sim = entities.NewSimulation(2020)

	// set up some initial entities.Sim.Companies
	for i := 0; i < 16; i++ {
		newCompany := economy.GenerateRandomCompany()
		entities.Sim.Companies = append(entities.Sim.Companies, newCompany)
		fmt.Printf("[ Econ ] %s (%s) founded!\n", newCompany.Name, newCompany.Industry)
	}

	ticker := time.NewTicker(1000 * time.Millisecond) // tick every second
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				entities.Sim.Tick()
				fmt.Println("")
				fmt.Printf("[ Date ] New simulation date is: %s\n", entities.Sim.Date)
				moveIn()
				findJobs()
				moveOut()

				// run entities.Sim.Market calculations every month
				diff := entities.Sim.Date.Sub(entities.Sim.Market.LastCalculation)
				if diff.Hours()/24 >= 28 {
					calculateEconomy()
				}
			}
		}
	}()

	jsonPtr := flag.Bool("json", false, "should output be in json?")
	// durationPtr := flag.Int("duration", 30, "how many seconds do we run the sim?")
	flag.Parse()

	// stop simulation after given duration
	ui.RunGame()
	ticker.Stop()
	done <- true

	printFinalState(*jsonPtr)
}

// people move in if there are free houses
func moveIn() {
	for i := 0; i < rand.Intn(freeHouses/4); i++ {
		h := people.CreateHousehold()
		freeHouses -= 1
		fmt.Printf("[ Move ] %s family has moved into a house, %d houses remain\n", h.FamilyName(), freeHouses)
		entities.Sim.People.MoveIn(h)
	}
}

// people move out if there are no jobs
func moveOut() {
	h := entities.Sim.People.Households
	// traverse in reverse order to avoid index shifting
	for i := len(h) - 1; i >= 0; i-- {
		if len(h[i].Members) > 0 && h[i].IsEligibleForMoveOut() {
			movedName := h[i].FamilyName()
			h = slices.Delete(h, i, i+1)
			freeHouses += 1
			fmt.Printf("[ Move ] %s family has moved out of the city, %d houses remain\n", movedName, freeHouses)
		}
	}
}

// assign unemployed people jobs
func findJobs() {
	h := entities.Sim.People.Households
	for i := 0; i < len(h); i++ {
		for j := 0; j < len(h[i].Members); j++ {
			if h[i].Members[j].IsEmployable() && !h[i].Members[j].IsEmployed() {
				companyId, remaining := getSuitableJob(h[i].Members[j])
				if companyId != 0 {
					h[i].Members[j].EmployerID = companyId
					fmt.Printf("[  Job ] %s %s has accepted a job as %s, %d jobs remain\n", h[i].Members[j].FirstName, h[i].Members[j].FamilyName, h[i].Members[j].Occupation, remaining)
				}
			}
		}
	}
}

func getSuitableJob(p entities.Person) (int, int) {
	remaining := 0
	companyId := 0
	for i := 0; i < len(entities.Sim.Companies); i++ {
		if entities.Sim.Companies[i].Industry == p.Industry {
			openings := entities.Sim.Companies[i].JobOpenings
			for j := 0; j < len(openings); j++ {
				if openings[p.CareerLevel] > 0 {
					entities.Sim.Companies[i].JobOpenings[p.CareerLevel] -= 1
					remaining = entities.Sim.Companies[i].GetNumberOfJobOpenings()
					companyId = entities.Sim.Companies[i].ID
				}
			}
		}
	}
	return companyId, remaining
}

func calculateEconomy() {
	// calculate impact of population growth on city economy
	population := entities.Sim.People.Population
	populationGrowth := entities.Sim.People.PopulationGrowthRate()

	entities.Sim.People.UpdatePopulationValues()
	entities.Sim.People.CalculateUnemployment()
	entities.Sim.Market.Unemployment = entities.Sim.People.UnemploymentRate()

	inflation := entities.Sim.Market.Inflation(populationGrowth)
	marketGrowth := entities.Sim.Market.MarketGrowth()
	newMarketValue := entities.Sim.Market.UpdateMarketValue(marketGrowth)

	fmt.Printf("[ Econ ] Town population is %d (±%.2f%%). Inflation: %.2f%%, Unemployment: %.2f%%, Market Value: %.2f (%.2f%%)\n", population, populationGrowth, inflation, entities.Sim.Market.Unemployment, newMarketValue, marketGrowth)

	if marketGrowth > 0 && rand.Intn(100) < 25 {
		newCompany := economy.GenerateRandomCompany()
		entities.Sim.Companies = append(entities.Sim.Companies, newCompany)
		fmt.Printf("[ Econ ] Growth! %s (%s) founded!\n", newCompany.Name, newCompany.Industry)
	}

	for k := 0; k < len(entities.Sim.Companies); k++ {
		entities.Sim.Companies[k].CalculateProfit()
		entities.Sim.Companies[k].DetermineJobOpenings()
	}
}

func printFinalState(printJson bool) {
	if printJson {
		cityDataJson, err := json.Marshal(entities.Sim.People)
		if err != nil {
			fmt.Println(err)
			return
		}
		compJson, err := json.Marshal(entities.Sim.Companies)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[ JSON ] Population: ", string(cityDataJson))
		fmt.Println("[ JSON ] Companies: ", string(compJson))
	} else {
		h := entities.Sim.People.Households
		for i := 0; i < len(h); i++ {
			for j := 0; j < len(h[i].Members); j++ {
				fmt.Println(h[i].Members[j].String())
			}
		}
		for k := 0; k < len(entities.Sim.Companies); k++ {
			fmt.Println(entities.Sim.Companies[k])
		}
	}

	fmt.Printf("[ Stat ] Total town population is %d (%.2f%% unemployment)\n", entities.Sim.People.Population, entities.Sim.Market.Unemployment)
}
