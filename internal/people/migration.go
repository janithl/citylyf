package people

import (
	"fmt"
	"math/rand/v2"

	"github.com/janithl/citylyf/internal/entities"
)

// Immigrate simulates inwards migration
func Immigrate() {
	if entities.Sim.Houses.GetFreeHouses() == 0 || rand.Float64() < 0.95 { // 5% change of moving in if there are free houses
		return
	}

	household := CreateHousehold()
	if houseID := household.FindHousing(); houseID > 0 {
		entities.Sim.People.Households[household.ID] = household
	} else {
		RemoveHousehold(household)
	}
}

// Emigrate simulates outwards migration
func Emigrate() {
	for _, household := range entities.Sim.People.Households {
		if household.Size() == 0 { // if a household is empty, remove it from the Sim and go to the next one
			delete(entities.Sim.People.Households, household.ID)
			continue
		}

		if household.IsEligibleForMoveOut() {
			movedName := household.FamilyName()
			houseID := household.HouseID
			RemoveHousehold(household)
			entities.Sim.Houses.MoveOut(houseID)
			fmt.Printf("[ Move ] %s family has moved out of house #%d and the city, %d houses remain\n", movedName, houseID, entities.Sim.Houses.GetFreeHouses())
		}
	}
}
