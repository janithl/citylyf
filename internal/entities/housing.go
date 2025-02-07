package entities

import (
	"fmt"
	"math/rand"
)

type Housing struct {
	Houses []House
}

func (h *Housing) MoveIn(bedrooms int) {
	for i := range h.Houses {
		if h.Houses[i].Free && h.Houses[i].Bedrooms >= bedrooms {
			h.Houses[i].Free = false
			return
		}
	}
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

type House struct {
	Free     bool
	Bedrooms int
}

func NewHousing(count int) *Housing {
	housing := &Housing{}
	for i := 0; i < count; i++ {
		housing.Houses = append(housing.Houses, House{
			Free:     true,
			Bedrooms: 1 + rand.Intn(3),
		})
	}
	fmt.Println(housing)
	return housing
}
