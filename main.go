package main

import (
	"citylyf/economy"
	"citylyf/entities"
	"citylyf/people"
	"fmt"
	"math/rand"
)

func getSuitableJob(companies []economy.Company, m economy.Market, p entities.Person) (int, int) {
	remaining := 0
	companyId := 0
	for i := 0; i < len(companies); i++ {
		if companies[i].Industry == p.Industry {

			openings := companies[i].JobOpenings
			for j := 0; j < len(openings); j++ {
				if openings[p.CareerLevel] > 0 {
					companies[i].JobOpenings[p.CareerLevel] -= 1
					remaining = companies[i].JobOpenings[p.CareerLevel]
					companyId = companies[i].ID
				}
			}
		}
	}
	return companyId, remaining
}

func main() {
	freeHouses := 100
	lastPopulation := 0
	unemployed := 0
	populationGrowth := 0.0

	market := economy.Market{
		InterestRate:           6.0,
		LastInflationRate:      0.0,
		Unemployment:           0.0,
		CorporateTax:           5.0,
		GovernmentSpending:     0.0,
		MonthsOfNegativeGrowth: 0,
	}

	var population []entities.Person
	var companies []economy.Company

	// set up some initial companies
	for i := 0; i < 5; i++ {
		newCompany := economy.GenerateRandomCompany(market)
		companies = append(companies, newCompany)
		fmt.Printf("%s (%s) founded!\n", newCompany.Name, newCompany.Industry)
	}

	for freeHouses > 0 {
		h := people.CreateHousehold()
		freeHouses -= 1
		fmt.Printf("%s family has moved into a house, %d houses remain\n", h.FamilyName(), freeHouses)

		for j := 0; j < len(h.Members); j++ {
			if h.Members[j].CareerLevel != entities.Unemployed {
				companyId, remaining := getSuitableJob(companies, market, h.Members[j])
				if companyId != 0 {
					h.Members[j].EmployerID = companyId
					fmt.Printf(">>> %s %s has accepted a job as %s, %d jobs remain\n", h.Members[j].FirstName, h.Members[j].FamilyName, h.Members[j].Occupation, remaining)
				} else {
					unemployed += 1
				}
			} else if h.Members[j].Age() > entities.AgeOfAdulthood {
				unemployed += 1
			}
		}

		population = append(population, h.Members...)

		// calculate impact of population growth on city economy
		populationGrowth = float64(len(population)-lastPopulation) / float64(lastPopulation)
		lastPopulation = len(population)
		market.Unemployment = 100 * float64(unemployed) / float64(lastPopulation)
		inflation := market.Inflation(populationGrowth)
		marketGrowth := market.MarketGrowth()
		fmt.Printf("Town population is %d (±%.2f%%). Inflation: %.2f%%, Market Growth: %.2f%%\n", len(population), populationGrowth, inflation, marketGrowth)

		if marketGrowth > 0 && rand.Intn(100) < 25 {
			newCompany := economy.GenerateRandomCompany(market)
			companies = append(companies, newCompany)
			fmt.Printf("Growth! %s (%s) founded!\n", newCompany.Name, newCompany.Industry)
		}

		for k := 0; k < len(companies); k++ {
			companies[k].CalculateProfit(market)
			companies[k].DetermineJobOpenings(market)
		}
	}

	fmt.Printf("Total town population is %d (%.2f%% unemployment)\n", len(population), market.Unemployment)
	for i := 0; i < len(population); i++ {
		fmt.Println(population[i].String())
	}
	for j := 0; j < len(companies); j++ {
		fmt.Println(companies[j])
	}
}
