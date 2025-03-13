package entities

import "math"

const (
	PopulationWeight    = 0.7
	JobAttractionWeight = 0.9
	ScalingConstant     = 0.4
)

type Trip struct {
	DestinationID, DailyTrips int
	Start, End                *Point
}

type Region struct {
	ID                            int
	Start                         Point
	Trips                         []*Trip
	Size, Population, Shops, Jobs int
}

func (r *Region) GetRegionalRoad() *Point {
	tiles := Sim.Geography.GetTiles()
	// Top and Bottom borders
	for x := r.Start.X; x < r.Start.X+r.Size; x++ {
		if Sim.Geography.BoundsCheck(x, r.Start.Y) && tiles[x][r.Start.Y].Road {
			return &Point{X: x, Y: r.Start.Y}
		}
		if Sim.Geography.BoundsCheck(x, r.Start.Y+r.Size-1) && tiles[x][r.Start.Y+r.Size-1].Road {
			return &Point{X: x, Y: r.Start.Y + r.Size - 1}
		}
	}
	// Left and Right borders
	for y := r.Start.Y; y < r.Start.Y+r.Size; y++ {
		if Sim.Geography.BoundsCheck(r.Start.X, y) && tiles[r.Start.X][y].Road {
			return &Point{X: r.Start.X, Y: y}
		}
		if Sim.Geography.BoundsCheck(r.Start.X+r.Size-1, y) && tiles[r.Start.X+r.Size-1][y].Road {
			return &Point{X: r.Start.X + r.Size - 1, Y: y}
		}
	}
	return nil
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
			r1road, r2road := r1.GetRegionalRoad(), r2.GetRegionalRoad()
			if r1road == nil || r2road == nil {
				continue
			}

			path := Sim.Geography.FindPath(r1road, r2road)
			if path == nil {
				continue
			}

			distance := float64(len(path))
			if distance < 1 {
				distance = 1 // Avoid division by zero
			}

			shoppingTrips := float64(r1.Population*r2.Shops) / math.Pow(distance, PopulationWeight)
			workTrips := float64(r1.Population*r2.Jobs) / math.Pow(distance, JobAttractionWeight)
			trips := int(math.Round(ScalingConstant * (workTrips + shoppingTrips)))
			if trips > 0 {
				r1.Trips = append(r1.Trips, &Trip{DestinationID: r2.ID, DailyTrips: trips, Start: r1road, End: r2road})
			}
		}
	}
}

func (r Regions) GetPopulationStats() ([][]int, int) {
	side := int(math.Sqrt(float64(len(r))))
	populationStats := make([][]int, side)
	maxPopulation := 0
	for x := 0; x < side; x++ {
		populationStats[x] = make([]int, side)
		for y := 0; y < side; y++ {
			rIndex := (x * side) + y
			if len(r) < rIndex {
				populationStats[x][y] = 0
				continue
			}
			regionPop := r[rIndex].Population
			populationStats[x][y] = regionPop
			if regionPop > maxPopulation {
				maxPopulation = regionPop
			}
		}
	}
	return populationStats, maxPopulation
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
