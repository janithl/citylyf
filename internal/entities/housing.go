package entities

import (
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
	ID               int
	Free             bool
	Bedrooms         int
	MonthlyRent      int
	LastRentRevision time.Time
}

type Housing struct {
	LastHouseID int
	Houses      []House
}

func (h *Housing) MoveIn(budget int, bedrooms int) int {
	for i := range h.Houses {
		if h.Houses[i].Free &&
			h.Houses[i].Bedrooms >= bedrooms &&
			h.Houses[i].MonthlyRent <= budget {
			h.Houses[i].Free = false
			h.Houses[i].LastRentRevision = Sim.Date // Lock in rents for 1 year
			return h.Houses[i].ID
		}
	}
	return 0
}

func (h *Housing) MoveOut() {
	for i := range h.Houses {
		if !h.Houses[i].Free {
			h.Houses[i].Free = true
			return
		}
	}
}

func (h *Housing) GetFreeHouses() int {
	freeCount := 0
	for _, house := range h.Houses {
		if house.Free {
			freeCount++
		}
	}
	return freeCount
}

func (h *Housing) GetHouse(id int) *House {
	for _, house := range h.Houses {
		if house.ID == id {
			return &house
		}
	}
	return nil
}

func (h *Housing) ReviseRents() {
	for i, house := range h.Houses {
		if Sim.Date.Sub(house.LastRentRevision).Hours() > HoursPerYear { // Revise rents every year
			h.Houses[i].MonthlyRent += int(float64(house.MonthlyRent) * Sim.Market.InterestRate() / 100)
			h.Houses[i].LastRentRevision = Sim.Date
		}
	}
}

func (h *Housing) AddHouse(x, y, bedrooms int) {
	if Sim.Geography.PlaceHouse(x, y, bedrooms < 4) { // house placed!
		h.LastHouseID++
		h.Houses = append(h.Houses, House{
			ID:               h.LastHouseID,
			Free:             true,
			Bedrooms:         bedrooms,
			MonthlyRent:      1200 + 200*(bedrooms-1),
			LastRentRevision: Sim.Date,
		})
	}
}
