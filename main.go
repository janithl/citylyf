package main

import (
	"citylyf/economy"
	"citylyf/people"
	"fmt"
)

func main() {
	freeHouses := 100
	lastPopulation := 0
	unemployed := 0
	populationGrowth := 0.0
	availableJobs := 60

	market := economy.Market{
		InterestRate:           6.0,
		LastInflationRate:      0.0,
		Unemployment:           0.0,
		CorporateTax:           5.0,
		GovernmentSpending:     0.0,
		MonthsOfNegativeGrowth: 0,
	}

	var population []people.Person

	for freeHouses > 0 && availableJobs > 0 {
		h := people.CreateHousehold()
		population = append(population, h.Members...)
		freeHouses -= 1
		fmt.Printf("%s family has moved into a house, %d houses remain\n", h.FamilyName(), freeHouses)

		for j := 0; j < len(h.Members); j++ {
			if h.Members[j].CareerLevel != people.Unemployed && availableJobs > 0 {
				availableJobs -= 1
				fmt.Printf("%s %s has accepted a job as %s, %d jobs remain\n", h.Members[j].FirstName, h.Members[j].FamilyName, h.Members[j].Occupation, availableJobs)
			} else if h.Members[j].Age() > people.AgeOfAdulthood {
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

		if marketGrowth > 0 {
			availableJobs = 1 + int(float64(availableJobs)*marketGrowth)
			fmt.Printf("Growth! %d jobs remain.\n", availableJobs)
		}

		// more than 6 months of negative growth means a recession, time for jobs to come down
		if market.MonthsOfNegativeGrowth > 6 {
			availableJobs = int(float64(availableJobs) * 0.66)
			fmt.Printf("Recession! %d jobs remain.\n", availableJobs)
		}
	}

	fmt.Printf("Total town population is %d\n", len(population))
	for i := 0; i < len(population); i++ {
		fmt.Println(population[i].String())
	}
}
