package entities

import (
	"maps"
	"math"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

type HouseType string

const (
	HouseSmall HouseType = "house-small"
	HouseLarge HouseType = "house-large"
)

type House struct {
	ID, HouseholdID, Bedrooms, MonthlyRent int
	HouseType                              HouseType
	Location                               *Point
	RoadDirection                          Direction
	LastRentRevision                       time.Time
}

type Housing map[int]*House

// GetIDs returns a sorted list of house IDs
func (h Housing) GetIDs() []int {
	IDs := []int{}
	for house := range maps.Values(h) {
		IDs = append(IDs, house.ID)
	}
	slices.Sort(IDs)
	return IDs
}

func (h Housing) MoveIn(householdID, budget, bedrooms int) int {
	for _, house := range h {
		if house.HouseholdID == 0 &&
			house.Bedrooms >= bedrooms &&
			house.MonthlyRent <= budget {
			house.HouseholdID = householdID
			house.LastRentRevision = Sim.Date // Lock in rents for 1 year
			return house.ID
		}
	}
	return 0
}

func (h Housing) MoveOut(houseID int) {
	house, exists := h[houseID]
	if exists {
		house.HouseholdID = 0
	}
}

func (h Housing) GetFreeHouses() int {
	freeCount := 0
	for _, house := range h {
		if house.HouseholdID == 0 {
			freeCount++
		}
	}
	return freeCount
}

func (h Housing) GetBaselineMonthlyRent(bedrooms int) int {
	if bedrooms < 2 {
		return 1200
	}

	return 1200 + 200*(bedrooms-1)
}

func (h Housing) GetAverageMonthlyRent() float64 {
	if len(h) == 0 {
		return 0.0
	}

	rent := 0.0
	for _, house := range h {
		rent += float64(house.MonthlyRent)
	}
	return rent / float64(len(h))
}

// GetCostOfLivingFactor returns a multiplier based on the change in average rent
func (h Housing) GetCostOfLivingFactor() float64 {
	costOfLivingFactor := h.GetAverageMonthlyRent() / float64(h.GetBaselineMonthlyRent(3))
	return math.Max(1.0, costOfLivingFactor) // has to be 1.0 or more
}

func (h Housing) VacancyRate() float64 {
	return float64(h.GetFreeHouses()) / float64(len(h))
}

func (h Housing) ReviseRents() {
	// Base adjustment using interest rate
	adjustmentFactor := Sim.Market.InterestRate() / 100

	// Adjust based on vacancy rate (high vacancy → reduce rent, low vacancy → normal increase)
	if h.VacancyRate() > 0.15 { // 15% vacancy threshold
		adjustmentFactor *= 0.5 // Reduce increase by 50%
	} else if h.VacancyRate() < 0.05 { // Low vacancy → demand is high
		adjustmentFactor *= 1.2 // Slightly boost increase
	}

	// Income & Inflation considerations: If income isn't rising, slow rent increases
	incomeFactor := Sim.People.AverageWageGrowthRate() - Sim.Market.InflationRate()
	if incomeFactor < 0 {
		adjustmentFactor *= 0.8 // Reduce rent increase if incomes are stagnating
	}

	// Apply the adjusted rent increase
	for _, house := range h {
		if Sim.Date.Sub(house.LastRentRevision).Hours() > HoursPerYear { // Revise rents every year
			currentRent := float64(house.MonthlyRent)
			rentChange := currentRent * adjustmentFactor
			rentChange = utils.Clamp(rentChange, -0.05*currentRent, 0.10*currentRent) // Clamp change between -5% and +10%

			house.MonthlyRent += int(rentChange)
			house.LastRentRevision = Sim.Date
		}
	}

	Sim.Market.History.AverageRent = utils.AddFifo(Sim.Market.History.AverageRent, h.GetAverageMonthlyRent(), 20)
}

// AverageRentGrowthRate returns the percentage growth rate of the AverageMonthlyRent
func (h Housing) AverageRentGrowthRate() float64 {
	if len(Sim.Market.History.AverageRent) == 0 {
		return 0.0
	}

	lastAverageRentValue := utils.GetLastValue(Sim.Market.History.AverageRent)
	if lastAverageRentValue == 0 {
		return 0.0
	}

	return 100.0 * (h.GetAverageMonthlyRent() - lastAverageRentValue) / lastAverageRentValue
}

func (h Housing) GetLocationHouse(x, y int) *House {
	for _, house := range h {
		if house.Location != nil && house.Location.X == x && house.Location.Y == y {
			return house
		}
	}
	return nil
}

func (h Housing) PlaceHousing() {
	if Sim.Market.HousingDemand < 0.05 && h.GetFreeHouses() > 3 { // low demand and enough free houses, no need to place more
		return
	}

	bedrooms := 2 + rand.IntN(3)
	site := Sim.Geography.GetPotentialSite(ResidentialUse)
	if site == nil { // no suitable sites
		return
	}

	Sim.Geography.tiles[site.X][site.Y].LandStatus = DevelopedStatus

	houseID := Sim.GetNextID()
	houseType := HouseSmall
	if bedrooms > 3 {
		houseType = HouseLarge
	}

	h[houseID] = &House{
		ID:               houseID,
		HouseholdID:      0,
		Bedrooms:         bedrooms,
		MonthlyRent:      int(h.GetCostOfLivingFactor() * float64(h.GetBaselineMonthlyRent(bedrooms))),
		HouseType:        houseType,
		Location:         site,
		RoadDirection:    Sim.Geography.getAccessRoad(site.X, site.Y),
		LastRentRevision: Sim.Date,
	}
}
