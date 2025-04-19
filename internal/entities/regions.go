package entities

import (
	"math"

	"github.com/janithl/citylyf/internal/utils"
)

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
	centre := Point{X: r.Start.X + r.Size/2, Y: r.Start.Y + r.Size/2}
	for d := range r.Size / 2 {
		for _, neighbour := range centre.GetNeighbours(d, false) {
			if Sim.Geography.BoundsCheck(neighbour.X, neighbour.Y) && tiles[neighbour.X][neighbour.Y].LandUse == TransportUse {
				return neighbour
			}
		}
	}
	return nil
}

func (r *Region) GetRegionalShops() []*Company {
	shops := []*Company{}
	tiles := Sim.Geography.GetTiles()
	for x := r.Start.X; x < r.Start.X+r.Size; x++ {
		for y := r.Start.Y; y < r.Start.Y+r.Size; y++ {
			if !Sim.Geography.BoundsCheck(x, y) {
				continue
			}

			switch {
			case tiles[x][y].LandUse == RetailUse:
				company := Sim.Companies.GetLocationCompany(x, y)
				if company != nil {
					shops = append(shops, company)
				}
			}
		}
	}
	return shops
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
				case tiles[x][y].LandUse == RetailUse:
					region.Shops += 1 // TODO: Check if shop is active
					company := Sim.Companies.GetLocationCompany(x, y)
					if company != nil {
						region.Jobs += company.GetNumberOfEmployees()
					}
				case tiles[x][y].LandUse == ResidentialUse:
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
	r.CalculateRegionalSales()
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

func (r Regions) CalculateRegionalSales() {
	for _, r1 := range r {
		if r1.Shops == 0 { // check if there are shops
			continue
		}

		population := float64(r1.Population)
		avgIncome := Sim.People.AverageMonthlyDisposableIncome()
		unemploymentRate := Sim.People.UnemploymentRate()
		consumerConfidence := utils.GetLastValue(Sim.Market.History.MarketSentiment)
		taxImpact := -Sim.Government.SalesTaxRate / 20
		inflationImpact := -math.Pow((Sim.Market.InflationRate()-5)/3, 2)

		// Base spending power calculation
		effectiveSpendingPower := float64(avgIncome) * (1 - unemploymentRate/1000) * (1 + consumerConfidence/10)
		effectiveSpendingPower = math.Max(0, effectiveSpendingPower) // Prevent negative values

		// Total demand within region
		regionRetailDemand := Sim.Market.RetailDemand * population * effectiveSpendingPower
		regionRetailDemand *= (1 + taxImpact + inflationImpact) // Adjust for macroeconomic factors

		// Add demand from outside the region
		externalRetailDemand := 0.0
		for _, r2 := range r {
			if r2.ID == r1.ID {
				continue
			}
			for _, trip := range r2.Trips {
				if trip.DestinationID != r1.ID {
					continue
				}

				tripSpendingPower := float64(avgIncome) * (1 - unemploymentRate/1000)
				// Assume a fraction of trips result in retail spending
				spendingTrips := float64(trip.DailyTrips) * 0.3 // 30% of trips include shopping
				externalRetailDemand += spendingTrips * tripSpendingPower
			}
		}

		// Get Total regional demand and distribute sales among shops
		totalRetailDemand := regionRetailDemand + externalRetailDemand
		avgSalesPerShop := totalRetailDemand / float64(r1.Shops)

		// Distribute sales among shops based on shop productivity
		for _, shop := range r1.GetRegionalShops() {
			shop.RetailSales = shop.GetProductivity() * avgSalesPerShop
		}
	}
}

func (r Regions) GetTotalTrips() int {
	trips := 0
	for _, region := range r {
		for _, trip := range region.Trips {
			trips += trip.DailyTrips
		}
	}
	return trips
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
