package entities

import (
	"fmt"
	"time"
)

type Household struct {
	Members    []Person  // Family members
	Savings    int       // Family savings
	MoveInDate time.Time // Day they moved in
}

func (h *Household) FamilyName() string {
	if len(h.Members) > 0 {
		return h.Members[0].FamilyName
	} else {
		return ""
	}
}

func (h *Household) AnnualIncome() int {
	income := 0
	for i := 0; i < len(h.Members); i++ {
		income += h.Members[i].AnnualIncome
	}
	return income
}

// eligible for move out if 1/4 years without income
func (h *Household) IsEligibleForMoveOut() bool {
	timeSinceMoveIn := Sim.Date.Sub(h.MoveInDate).Hours() / HoursPerYear
	noIncome := true
	for i := 0; i < len(h.Members); i++ {
		if h.Members[i].IsEmployed() {
			noIncome = false
		}
	}
	return noIncome && timeSinceMoveIn > 0.25
}

func (h *Household) GetStats() string {
	return fmt.Sprintf("%s family (%d members) - Moved in %s", h.FamilyName(), len(h.Members), h.MoveInDate.Format("2006-01-02"))
}
