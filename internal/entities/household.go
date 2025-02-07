package entities

import (
	"fmt"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

type Household struct {
	Members    []Person  // Family members
	Savings    int       // Family savings
	LastPayDay time.Time // Last time payments were calculated
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

// calculate monthly pay
// TODO extend to calculate expenses as well
func (h *Household) CalculateMonthlyPay() {
	daysSinceLastPay := Sim.Date.Sub(h.LastPayDay).Hours() / HoursPerDay
	pay := 0.0
	for i := range h.Members {
		if h.Members[i].IsEmployed() {
			memberPay := float64(h.Members[i].AnnualIncome) * daysSinceLastPay / DaysPerYear
			h.Members[i].Savings += int(memberPay)
			pay += memberPay
		}
	}
	h.Savings += int(pay)
	h.LastPayDay = Sim.Date
}

func (h *Household) GetStats() string {
	return fmt.Sprintf("%-24s %d Members   %s   %s", h.FamilyName()+" family", len(h.Members),
		"Moved in "+h.MoveInDate.Format("2006-01-02"), utils.FormatCurrency(float64(h.Savings), "$"))
}
