package entities

import (
	"maps"
	"slices"
	"time"
)

type HouseType string

const (
	NonHouse        HouseType = ""
	HouseSmallX     HouseType = "house-small-x"
	HouseSmallXBack HouseType = "house-small-x-back"
	HouseSmallY     HouseType = "house-small-y"
	HouseSmallYBack HouseType = "house-small-y-back"
	HouseLargeX     HouseType = "house-large-x"
	HouseLargeXBack HouseType = "house-large-x-back"
	HouseLargeY     HouseType = "house-large-y"
	HouseLargeYBack HouseType = "house-large-y-back"
)

type House struct {
	ID, HouseholdID, Bedrooms, MonthlyRent int
	LastRentRevision                       time.Time
}

type Housing map[int]*House

// GetIDs returns a sorted list of house IDs
func (h Housing) GetIDs() []int {
	IDs := []int{}
	for house := range maps.Values(h) {
		IDs = append(IDs, house.ID)
	}
	slices.Sort(IDs)
	return IDs
}

func (h Housing) MoveIn(householdID, budget, bedrooms int) int {
	for _, id := range h.GetIDs() {
		if h[id].HouseholdID == 0 &&
			h[id].Bedrooms >= bedrooms &&
			h[id].MonthlyRent <= budget {
			h[id].HouseholdID = householdID
			h[id].LastRentRevision = Sim.Date // Lock in rents for 1 year
			return h[id].ID
		}
	}
	return 0
}

func (h Housing) MoveOut(houseID int) {
	house, exists := h[houseID]
	if exists {
		house.HouseholdID = 0
	}
}

func (h Housing) GetFreeHouses() int {
	freeCount := 0
	for _, id := range h.GetIDs() {
		if h[id].HouseholdID == 0 {
			freeCount++
		}
	}
	return freeCount
}

func (h Housing) ReviseRents() {
	for _, id := range h.GetIDs() {
		if Sim.Date.Sub(h[id].LastRentRevision).Hours() > HoursPerYear { // Revise rents every year
			h[id].MonthlyRent += int(float64(h[id].MonthlyRent) * Sim.Market.InterestRate() / 100)
			h[id].LastRentRevision = Sim.Date
		}
	}
}

func (h Housing) AddHouse(x, y, bedrooms int) {
	if Sim.Geography.PlaceHouse(x, y, bedrooms < 4) { // house placed!
		houseID := Sim.GetNextID()
		h[houseID] = &House{
			ID:               houseID,
			HouseholdID:      0,
			Bedrooms:         bedrooms,
			MonthlyRent:      1200 + 200*(bedrooms-1),
			LastRentRevision: Sim.Date,
		}
	}
}
