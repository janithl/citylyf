package people

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/janithl/citylyf/entities"
)

// Migration handles people moving in and out of town
type Migration struct {
	FreeHouses int // we'll keep this here for now until we move it
}

// people move in if there are free houses
func (m *Migration) MoveIn() {
	for i := 0; i < rand.Intn(1+(m.FreeHouses/4)); i++ {
		h := CreateHousehold()
		m.FreeHouses--
		fmt.Printf("[ Move ] %s family has moved into a house, %d houses remain\n", h.FamilyName(), m.FreeHouses)
		entities.Sim.People.MoveIn(h)
	}
}

// people move out if there are no jobs
func (m *Migration) MoveOut() {
	h := entities.Sim.People.Households
	// traverse in reverse order to avoid index shifting
	for i := len(h) - 1; i >= 0; i-- {
		if len(h[i].Members) > 0 && h[i].IsEligibleForMoveOut() {
			movedName := h[i].FamilyName()
			h = slices.Delete(h, i, i+1)
			m.FreeHouses++
			fmt.Printf("[ Move ] %s family has moved out of the city, %d houses remain\n", movedName, m.FreeHouses)
		}
	}
	entities.Sim.People.Households = h
}
