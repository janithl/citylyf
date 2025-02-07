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

// TODO
// Household Budgeting - think about rent/mortgage expenses + taxes + savings interest etc
// Housing market - rent, no. of bedrooms etc.
// People should marry, have babies, move out out the house, die etc.
// Yearly budget - once a year, we calculate how much income tax and company tax has to be paid
func main() {
	entities.Sim = entities.NewSimulation(2020, 10+rand.Intn(10), 100000)
	employment := economy.Employment{}

	// set up some initial entities.Sim.Companies
	for i := 0; i < 4+rand.Intn(4); i++ {
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
				entities.Sim.People.MoveIn(people.CreateHousehold)
				employment.AssignJobs()
				entities.Sim.People.MoveOut()

				// run entities.Sim.Market calculations every month
				diff := entities.Sim.Date.Sub(entities.Sim.Market.LastCalculation)
				if diff.Hours()/24 >= 28 {
					calculateEconomy()
					entities.Sim.Government.CollectTaxes()
				}
			}
		}
	}()

	jsonPtr := flag.Bool("json", false, "should output be in json?")
	flag.Parse()

	// run simulation until UI is closed
	ui.RunGame()
	ticker.Stop()
	done <- true

	if *jsonPtr {
		printFinalState()
	}
}

func calculateEconomy() {
	// calculate impact of population growth on city economy
	populationGrowth := entities.Sim.People.PopulationGrowthRate()

	entities.Sim.People.UpdatePopulationValues()
	entities.Sim.People.CalculateUnemployment()
	entities.Sim.Market.Unemployment = entities.Sim.People.UnemploymentRate()

	entities.Sim.Market.Inflation(populationGrowth)
	marketGrowth := entities.Sim.Market.MarketGrowth()
	entities.Sim.Market.UpdateMarketValue(marketGrowth)

	fmt.Printf("[ Econ ] %s\n", entities.Sim.GetStats())

	if marketGrowth > 0 && rand.Intn(100) < 5 { // 5% chance of a company being formed during the good times
		newCompany := economy.GenerateRandomCompany()
		entities.Sim.Companies = append(entities.Sim.Companies, newCompany)
		fmt.Printf("[ Econ ] Growth! %s (%s) founded!\n", newCompany.Name, newCompany.Industry)
	}

	totalProfits := 0.0
	for k := 0; k < len(entities.Sim.Companies); k++ {
		totalProfits += entities.Sim.Companies[k].CalculateProfit()
		entities.Sim.Companies[k].DetermineJobOpenings()
	}
	entities.Sim.Market.ReportCompanyProfits(totalProfits)

	// calculate monthly pay
	for i := range entities.Sim.People.Households {
		entities.Sim.People.Households[i].CalculateMonthlyPay()
	}

	// do interest calcuations
	monthlyInterestRate := entities.Sim.Market.InterestRate * 28 / 36500
	entities.Sim.Government.Reserves += int(float64(entities.Sim.Government.Reserves) * monthlyInterestRate)
	for i := range entities.Sim.People.Households {
		entities.Sim.People.Households[i].Savings += int(float64(entities.Sim.People.Households[i].Savings) * monthlyInterestRate)
	}
}

func printFinalState() {
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
}
