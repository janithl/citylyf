package main

import (
	"citylyf/people"
	"fmt"
)

func main() {
	for i := 0; i < 100; i++ {
		h := people.CreateHousehold()
		fmt.Printf("[The %s Family | $%8d/yearly | $%8d ]:\n", h.FamilyName(), h.AnnualIncome(), h.Wealth())
		for j := 0; j < len(h.Members); j++ {
			fmt.Printf(" |-> %s\n", h.Members[j].String())
		}
	}
}
