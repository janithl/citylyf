package internal_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/janithl/citylyf/internal"
	"github.com/janithl/citylyf/internal/entities"
)

func BenchmarkSimRunner(b *testing.B) {
	// Set up the simulation
	simRunner := &internal.SimRunner{}
	simRunner.NewGame(nil)
	entities.Sim.SimulationSpeed = entities.Slow

	// add some roads and residential zones
	simSize := entities.Sim.Geography.Size
	for i := 0; i < 16; i++ {
		x, y := rand.IntN(simSize), rand.IntN(simSize)
		entities.PlaceRoad(entities.Point{X: x - 1, Y: y}, entities.Point{X: x + 1, Y: y}, entities.Asphalt)
		use := entities.ResidentialUse
		if i >= 12 {
			use = entities.RetailUse
		}
		entities.Sim.Geography.PlaceLandUse(entities.Point{X: x - 2, Y: y - 1}, entities.Point{X: x + 2, Y: y + 1}, use)
	}

	for b.Loop() {
		entities.Sim.Tick(simRunner.GameTick)
		entities.Sim.SendStats()
	}

	for _, household := range entities.Sim.People.Households {
		fmt.Println(household.GetStats())
		fmt.Println(household.GetMemberStats())
	}
}
