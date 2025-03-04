package entities

import (
	"maps"
	"slices"
	"time"
)

type HouseType string

const (
	HouseSmall HouseType = "house-small"
	HouseLarge HouseType = "house-large"
)

type House struct {
	ID, HouseholdID, Bedrooms, MonthlyRent int
	HouseType                              HouseType
	Location                               Point
	RoadDirection                          Direction
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

func (h Housing) GetLocationHouse(x, y int) *House {
	for house := range maps.Values(h) {
		if house.Location.X == x && house.Location.Y == y {
			return house
		}
	}
	return nil
}

func (h Housing) AddHouse(x, y, bedrooms int) {
	if roadDir := Sim.Geography.getAccessRoad(x, y); roadDir != "" { // house placeable!
		Sim.Geography.placeHouse(x, y)
		houseID := Sim.GetNextID()
		houseType := HouseSmall
		if bedrooms > 3 {
			houseType = HouseLarge
		}

		h[houseID] = &House{
			ID:               houseID,
			HouseholdID:      0,
			Bedrooms:         bedrooms,
			MonthlyRent:      1200 + 200*(bedrooms-1),
			HouseType:        houseType,
			Location:         Point{x, y},
			RoadDirection:    roadDir,
			LastRentRevision: Sim.Date,
		}
	}
}
