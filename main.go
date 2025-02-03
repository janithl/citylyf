package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/janithl/citylyf/economy"
	"github.com/janithl/citylyf/entities"
	"github.com/janithl/citylyf/people"
)

// TODO
// Move people out of the city when they can't find work after a certain time - think about rent/mortgage expenses
var population []entities.Person
var companies []economy.Company
var freeHouses int
var lastPopulation int
var unemployed int
var populationGrowth float64
var market economy.Market

func main() {
	entities.CitySimulation = entities.NewSimulation(2020)

	freeHouses = 100
	lastPopulation = 0
	populationGrowth = 0.0
	market = economy.Market{
		InterestRate:           5.0,
		LastInflationRate:      0.0,
		Unemployment:           0.0,
		CorporateTax:           3.0,
		GovernmentSpending:     10.0,
		MonthsOfNegativeGrowth: 0,
		LastCalculation:        entities.CitySimulation.Date,
	}

	// set up some initial companies
	for i := 0; i < 16; i++ {
		newCompany := economy.GenerateRandomCompany(market)
		companies = append(companies, newCompany)
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
				entities.CitySimulation.Tick()
				fmt.Println("")
				fmt.Printf("[ Date ] New simulation date is: %s\n", entities.CitySimulation.Date)
				moveIn()
				findJobs()

				// run market calculations every month
				diff := entities.CitySimulation.Date.Sub(market.LastCalculation)
				if diff.Hours()/24 >= 28 {
					calculateEconomy()
				}
			}
		}
	}()

	jsonPtr := flag.Bool("json", false, "should output be in json?")
	durationPtr := flag.Int("duration", 30, "how many seconds do we run the sim?")
	flag.Parse()

	// stop simulation after given duration
	time.Sleep(time.Duration(*durationPtr) * time.Second)
	ticker.Stop()
	done <- true

	printFinalState(*jsonPtr)
}

func moveIn() {
	for i := 0; i < rand.Intn(freeHouses/2); i++ {
		h := people.CreateHousehold()
		freeHouses -= 1
		fmt.Printf("[ Move ] %s family has moved into a house, %d houses remain\n", h.FamilyName(), freeHouses)
		population = append(population, h.Members...)
	}
}

func findJobs() {
	// assign unemployed people jobs
	unemployed = 0
	for j := 0; j < len(population); j++ {
		if population[j].EmployerID == 0 && population[j].CareerLevel != entities.Unemployed {
			companyId, remaining := getSuitableJob(companies, market, population[j])
			if companyId != 0 {
				population[j].EmployerID = companyId
				fmt.Printf("[  Job ] %s %s has accepted a job as %s, %d jobs remain\n", population[j].FirstName, population[j].FamilyName, population[j].Occupation, remaining)
			} else {
				unemployed += 1
			}
		} else if population[j].EmployerID == 0 && population[j].Age() > entities.AgeOfAdulthood {
			unemployed += 1
		}
	}
}

func getSuitableJob(companies []economy.Company, m economy.Market, p entities.Person) (int, int) {
	remaining := 0
	companyId := 0
	for i := 0; i < len(companies); i++ {
		if companies[i].Industry == p.Industry {
			openings := companies[i].JobOpenings
			for j := 0; j < len(openings); j++ {
				if openings[p.CareerLevel] > 0 {
					companies[i].JobOpenings[p.CareerLevel] -= 1
					remaining = companies[i].GetNumberOfJobOpenings()
					companyId = companies[i].ID
				}
			}
		}
	}
	return companyId, remaining
}

func calculateEconomy() {
	// calculate impact of population growth on city economy
	populationGrowth = float64(len(population)-lastPopulation) / float64(lastPopulation)
	lastPopulation = len(population)
	market.Unemployment = 100 * float64(unemployed) / float64(lastPopulation)
	inflation := market.Inflation(populationGrowth)
	marketGrowth := market.MarketGrowth()
	market.LastCalculation = entities.CitySimulation.Date // update last calculation time

	fmt.Printf("[ Econ ] Town population is %d (±%.2f%%). Inflation: %.2f%%, Unemployment: %.2f%%, Market Growth: %.2f%%\n", len(population), populationGrowth, inflation, market.Unemployment, marketGrowth)

	if marketGrowth > 0 && rand.Intn(100) < 50 {
		newCompany := economy.GenerateRandomCompany(market)
		companies = append(companies, newCompany)
		fmt.Printf("[ Econ ] Growth! %s (%s) founded!\n", newCompany.Name, newCompany.Industry)
	}

	for k := 0; k < len(companies); k++ {
		companies[k].CalculateProfit(market)
		companies[k].DetermineJobOpenings(market)
	}
}

func printFinalState(printJson bool) {
	if printJson {
		popJson, err := json.Marshal(population)
		if err != nil {
			fmt.Println(err)
			return
		}
		compJson, err := json.Marshal(companies)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[ JSON ] Population: ", string(popJson))
		fmt.Println("[ JSON ] Companies: ", string(compJson))
	} else {
		for i := 0; i < len(population); i++ {
			fmt.Println(population[i].String())
		}
		for j := 0; j < len(companies); j++ {
			fmt.Println(companies[j])
		}
	}

	fmt.Printf("[ Stat ] Total town population is %d (%.2f%% unemployment)\n", len(population), market.Unemployment)
}
