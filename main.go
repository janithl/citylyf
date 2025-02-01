package main

import (
	"citylyf/people"
	"fmt"
)

func main() {
	freeHouses := 100
	availableJobs := 30
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
			}
		}
	}

	fmt.Printf("Total town population is %d\n", len(population))
}
