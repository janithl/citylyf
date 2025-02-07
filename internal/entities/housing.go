package entities

import (
	"math/rand"
)

type Housing struct {
	Houses []House
}

func (h *Housing) MoveIn(budget int, bedrooms int) int {
	for i := range h.Houses {
		if h.Houses[i].Free &&
			h.Houses[i].Bedrooms >= bedrooms &&
			h.Houses[i].MonthlyRent <= budget {
			h.Houses[i].Free = false
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

type House struct {
	ID          int
	Free        bool
	Bedrooms    int
	MonthlyRent int
}

func NewHousing(count int) *Housing {
	housing := &Housing{}
	for i := 0; i < count; i++ {
		bedrooms := 1 + rand.Intn(3)
		housing.Houses = append(housing.Houses, House{
			ID:          100 + i,
			Free:        true,
			Bedrooms:    bedrooms,
			MonthlyRent: 1200 + 200*(bedrooms-1),
		})
	}
	return housing
}
