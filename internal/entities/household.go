package entities

import (
	"fmt"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

type Household struct {
	ID, HouseID int       // Family ID and the ID of the house they live at
	MemberIDs   []int     // Family member IDs
	Savings     int       // Family savings
	LastPayDay  time.Time // Last time payments were calculated
	MoveInDate  time.Time // Day they moved in
}

func (h *Household) Size() int {
	return len(h.MemberIDs)
}

func (h *Household) FamilyName() string {
	if h.Size() > 0 {
		p := Sim.People.GetPerson(h.MemberIDs[0])
		if p != nil {
			return p.FamilyName
		}
	}
	return ""
}

func (h *Household) AnnualIncome() int {
	income := 0
	for _, memberID := range h.MemberIDs {
		p := Sim.People.GetPerson(memberID)
		if p != nil {
			income += p.AnnualIncome
		}
	}
	return income
}

// eligible for move out if 1/4 years without income
func (h *Household) IsEligibleForMoveOut() bool {
	timeSinceMoveIn := Sim.Date.Sub(h.MoveInDate).Hours() / HoursPerYear
	noIncome := true
	for _, memberID := range h.MemberIDs {
		p := Sim.People.GetPerson(memberID)
		if p != nil && (p.IsEmployed() || (p.CareerLevel == Retired && h.Savings > 0)) {
			noIncome = false
		}
	}
	return noIncome && timeSinceMoveIn > 0.25
}

// calculate monthly budget
func (h *Household) CalculateMonthlyBudget(addPayToPayroll func(companyID int, payAmount float64)) {
	daysSinceLastPay := Sim.Date.Sub(h.LastPayDay).Hours() / HoursPerDay
	pay := 0.0
	for _, memberID := range h.MemberIDs {
		p := Sim.People.GetPerson(memberID)
		if p != nil && p.IsEmployed() {
			memberPay := float64(p.AnnualIncome) * daysSinceLastPay / DaysPerYear
			p.Savings += int(memberPay)
			addPayToPayroll(p.EmployerID, memberPay) // deduct from company
			pay += memberPay
		}
	}
	house, exists := Sim.Houses[h.HouseID]
	if exists {
		h.Savings += int(pay) - house.MonthlyRent
		h.LastPayDay = Sim.Date
	}
}

func (h *Household) GetID() int {
	return h.ID
}

func (h *Household) GetStats() string {
	return fmt.Sprintf("%-24s %d Membs   %s   %s", h.FamilyName()+" family", h.Size(),
		"Moved in "+h.MoveInDate.Format("2006-01-02"), utils.FormatCurrency(float64(h.Savings), "$"))
}

func (h *Household) GetMemberStats() string {
	stats := ""
	for _, memberID := range h.MemberIDs {
		p := Sim.People.GetPerson(memberID)
		if p != nil {
			stats += p.String() + "\n"
		}
	}
	return stats
}
