package main

import (
	"citylyf/economy"
	"citylyf/entities"
	"citylyf/people"
	"fmt"
	"math/rand"
)

func getSuitableJob(companies []economy.Company, m economy.Market, p entities.Person) int {
	remaining := -1
	for i := 0; i < len(companies); i++ {
		if companies[i].Industry == p.Industry {

			openings := companies[i].JobOpenings
			for j := 0; j < len(openings); j++ {
				if openings[p.CareerLevel] > 0 {
					companies[i].JobOpenings[p.CareerLevel] -= 1
					remaining = companies[i].JobOpenings[p.CareerLevel]
				}
			}
		}
	}
	return remaining
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

	for freeHouses > 0 {
		h := people.CreateHousehold()
		population = append(population, h.Members...)
		freeHouses -= 1
		fmt.Printf("%s family has moved into a house, %d houses remain\n", h.FamilyName(), freeHouses)

		for j := 0; j < len(h.Members); j++ {
			if h.Members[j].CareerLevel != entities.Unemployed {
				remaining := getSuitableJob(companies, market, h.Members[j])
				if remaining > -1 {
					fmt.Printf(">>> %s %s has accepted a job as %s, %d jobs remain\n", h.Members[j].FirstName, h.Members[j].FamilyName, h.Members[j].Occupation, remaining)
				}
			} else if h.Members[j].Age() > entities.AgeOfAdulthood {
				unemployed += 1
			}
		}

		// calculate impact of population growth on city economy
		populationGrowth = float64(len(population)-lastPopulation) / float64(lastPopulation)
		lastPopulation = len(population)
		market.Unemployment = 100 * float64(unemployed) / float64(lastPopulation)
		inflation := market.Inflation(populationGrowth)
		marketGrowth := market.MarketGrowth()
		fmt.Printf("Town population is %d (±%.2f%%). Inflation: %.2f%%, Market Growth: %.2f%%\n", len(population), populationGrowth, inflation, marketGrowth)

		if marketGrowth > 0 && rand.Intn(100) < 25 {
			newCompany := economy.Company{
				Name:         fmt.Sprintf("Such a Place %d", len(companies)),
				Industry:     entities.Software,
				LastRevenue:  100_000,
				LastExpenses: 75_000,
				LastProfit:   25_000,
				JobOpenings:  make(map[entities.CareerLevel]int),
			}
			companies = append(companies, newCompany)

			fmt.Printf("Growth! 1 company started %s\n", newCompany.Name)
		}

		// more than 6 months of negative growth means a recession, time for jobs to come down
		// if market.MonthsOfNegativeGrowth > 6 {
		// 	availableJobs = int(float64(availableJobs) * 0.66)
		// 	fmt.Printf("Recession! %d jobs remain.\n", availableJobs)
		// }
		for k := 0; k < len(companies); k++ {
			companies[k].CalculateProfit(market)
			companies[k].DetermineJobOpenings(market)
		}
	}

	fmt.Printf("Total town population is %d\n", len(population))
	for i := 0; i < len(population); i++ {
		fmt.Println(population[i].String())
	}
	for j := 0; j < len(companies); j++ {
		fmt.Println(companies[j])
	}
}
