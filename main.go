package main

import (
	"citylyf/economy"
	"citylyf/people"
	"fmt"
)

func main() {
	freeHouses := 100
	availableJobs := 30
	lastPopulation := 0
	unemployed := 0
	populationGrowth := 0.0

	market := economy.Market{
		InterestRate:       5.0,
		LastInflationRate:  0.0,
		Unemployment:       0.0,
		CorporateTax:       5.0,
		GovernmentSpending: 0.0,
	}

	var population []people.Person

	for freeHouses > 0 && availableJobs > 0 {
		h := people.CreateHousehold()
		population = append(population, h.Members...)
		freeHouses -= 1
		fmt.Printf("%s family has moved into a house, %d houses remain\n", h.FamilyName(), freeHouses)

		for j := 0; j < len(h.Members); j++ {
			if h.Members[j].CareerLevel != people.Unemployed {
				availableJobs -= 1
				fmt.Printf("%s %s has accepted a job as %s\n", h.Members[j].FirstName, h.Members[j].FamilyName, h.Members[j].Occupation)
			} else {
				unemployed += 1
			}
		}

		// calculate impact of population growth on city economy
		populationGrowth = float64(len(population)-lastPopulation) / float64(lastPopulation)
		lastPopulation = len(population)
		market.Unemployment = 100 * float64(unemployed) / float64(lastPopulation)
		inflation := market.Inflation(populationGrowth)
		market.LastInflationRate = inflation
		fmt.Printf("Town population is %d (±%.2f%%). Inflation is at %.2f%%\n", len(population), populationGrowth, inflation)
	}

	fmt.Printf("Total town population is %d\n", len(population))
	for i := 0; i < len(population); i++ {
		fmt.Println(population[i].String())
	}
}
