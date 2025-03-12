package entities

const (
	PopulationWeight    = 0.7
	JobAttractionWeight = 0.9
)

type Trip struct {
	DestinationID, MonthlyCount int
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
