package entities

import "math"

const (
	PopulationWeight    = 0.7
	JobAttractionWeight = 0.9
	ScalingConstant     = 0.15
)

type Trip struct {
	DestinationID, DailyTrips int
}

type Region struct {
	ID                            int
	Start                         Point
	Trips                         []*Trip
	Size, Population, Shops, Jobs int
}

type Regions []*Region

func (r Regions) CalculateRegionalStats() {
	tiles := Sim.Geography.GetTiles()
	for _, region := range r {
		region.Shops = 0
		region.Jobs = 0
		region.Population = 0
		for x := region.Start.X; x < region.Start.X+region.Size; x++ {
			for y := region.Start.Y; y < region.Start.Y+region.Size; y++ {
				if !Sim.Geography.BoundsCheck(x, y) {
					continue
				}

				switch {
				case tiles[x][y].Shop:
					region.Shops += 1
					company := Sim.Companies.GetLocationCompany(x, y)
					if company != nil {
						region.Jobs += company.GetNumberOfEmployees()
					}
				case tiles[x][y].House:
					house := Sim.Houses.GetLocationHouse(x, y)
					if house != nil && house.HouseholdID != 0 {
						household, exists := Sim.People.Households[house.HouseholdID]
						if exists {
							region.Population += household.Size()
						}
					}
				}
			}
		}
	}
	r.CalculateRegionalTraffic()
}

func (r Regions) CalculateRegionalTraffic() {
	for _, r1 := range r {
		r1.Trips = []*Trip{}
		for _, r2 := range r {
			if r1.ID == r2.ID {
				continue // Ignore same-region trips
			}

			distance := math.Sqrt(math.Pow(float64(r1.Start.X)-float64(r2.Start.X), 2) +
				math.Pow(float64(r1.Start.Y)-float64(r2.Start.Y), 2)) // TODO: Improve this

			if distance < 1 {
				distance = 1 // Avoid division by zero
			}

			shoppingTrips := float64(r1.Population*r2.Shops) / math.Pow(distance, PopulationWeight)
			workTrips := float64(r1.Population*r2.Jobs) / math.Pow(distance, JobAttractionWeight)
			trips := int(math.Round(ScalingConstant * (workTrips + shoppingTrips)))
			if trips > 0 {
				r1.Trips = append(r1.Trips, &Trip{DestinationID: r2.ID, DailyTrips: trips})
			}
		}
	}
}

func NewRegions(mapSize, regionSize int) []*Region {
	regions := []*Region{}

	side := mapSize / regionSize
	for x := 0; x < side; x++ {
		for y := 0; y < side; y++ {
			region := &Region{
				ID:         len(regions) + 1,
				Start:      Point{X: x * side, Y: y * side},
				Trips:      []*Trip{},
				Size:       regionSize,
				Population: 0,
				Shops:      0,
				Jobs:       0,
			}
			regions = append(regions, region)
		}
	}

	return regions
}
