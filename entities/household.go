package entities

import "time"

type Household struct {
	Members    []Person  // Family members
	Savings    int       // Family savings
	MoveInDate time.Time // Day they moved in
}

func (h *Household) FamilyName() string {
	return h.Members[0].FamilyName
}

func (h *Household) AnnualIncome() int {
	income := 0
	for i := 0; i < len(h.Members); i++ {
		income += h.Members[i].AnnualIncome
	}
	return income
}
