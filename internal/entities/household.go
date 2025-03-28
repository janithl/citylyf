package entities

import (
	"fmt"
	"slices"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

type Household struct {
	ID, HouseID       int       // Family ID and the ID of the house they live at
	MemberIDs         []int     // Family member IDs
	Savings           int       // Family savings
	LastMonthExpenses int       // total expenses last month
	LastPayDay        time.Time // Last time payments were calculated
	MoveInDate        time.Time // Day they moved in
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

func (h *Household) GetMembers() []*Person {
	members := []*Person{}
	for _, memberID := range h.MemberIDs {
		p := Sim.People.GetPerson(memberID)
		if p != nil {
			members = append(members, p)
		}
	}
	return members
}

func (h *Household) AddMember(personID int, savings int) {
	h.MemberIDs = append(h.MemberIDs, personID)
	h.Savings += savings
}

// if potential = true, this returns the ideal annual income if all employeable people are employed
func (h *Household) AnnualIncome(potential bool) int {
	income := 0
	for _, memberID := range h.MemberIDs {
		p := Sim.People.GetPerson(memberID)
		if p != nil {
			if potential && p.IsEmployable() {
				income += p.AnnualIncome
			} else {
				income += p.CurrentIncome()
			}
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
		if p != nil {
			memberPay := float64(p.CurrentIncome()) * daysSinceLastPay / DaysPerYear
			p.Savings += int(memberPay)
			addPayToPayroll(p.EmployerID, memberPay) // deduct from company
			pay += memberPay
		}
	}
	house, exists := Sim.Houses[h.HouseID]
	if exists {
		expenses := house.MonthlyRent // TODO: Expand this
		h.Savings += int(pay) - expenses
		h.LastMonthExpenses = expenses
		h.LastPayDay = Sim.Date
	}
}

func (h *Household) GetID() int {
	return h.ID
}

func (h *Household) GetStats() string {
	return fmt.Sprintf("%-30s %02d/%02d   %s   %s", h.FamilyName()+" family", h.Size(), h.GetEmployedCount(),
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

// GetEmployedCount returns the number of employed members of the household
func (h *Household) GetEmployedCount() int {
	employed := 0
	for _, memberID := range h.MemberIDs {
		p := Sim.People.GetPerson(memberID)
		if p != nil && p.IsEmployed() {
			employed++
		}
	}
	return employed
}

// GetAdultCount returns the number of adult members of the household
func (h *Household) GetAdultCount() int {
	adults := 0
	for _, memberID := range h.MemberIDs {
		p := Sim.People.GetPerson(memberID)
		if p != nil && p.Age() >= AgeOfAdulthood {
			adults++
		}
	}
	return adults
}

// IsMember returns true if a given person id is a member of the household
func (h *Household) IsMember(personID int) bool {
	for _, memberID := range h.MemberIDs {
		if memberID == personID {
			return true
		}
	}
	return false
}

// RemoveMember returns removes the given person from the household
func (h *Household) RemoveMember(person *Person) {
	if !h.IsMember(person.ID) {
		return
	}

	h.MemberIDs = slices.DeleteFunc(h.MemberIDs, func(id int) bool {
		return id == person.ID
	})
	h.Savings -= person.Savings
}

// FindHousing assigns a house to a househld
func (h *Household) FindHousing() int {
	monthlyRentBudget := float64(h.AnnualIncome(true)) / (4 * 12)          // 25% of (potential) yearly income towards rent / 12
	houseID := Sim.Houses.MoveIn(h.ID, int(monthlyRentBudget), h.Size()/2) // everyone gets to share a bedroom
	if houseID > 0 {
		h.HouseID = houseID
		fmt.Printf("[ Move ] %s family has moved into house #%d, %d houses remain\n", h.FamilyName(), houseID, Sim.Houses.GetFreeHouses())
	}

	return houseID
}
