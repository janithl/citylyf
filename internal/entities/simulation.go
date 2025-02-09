package entities

import (
	"fmt"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

// SimulationSpeed defines the number of days at which the sim moves
type SimulationSpeed int

const (
	Pause SimulationSpeed = 0
	Slow  SimulationSpeed = 3
	Mid   SimulationSpeed = 7
	Fast  SimulationSpeed = 28
)

type Simulation struct {
	SimulationSpeed SimulationSpeed
	Date            time.Time
	Government      Government
	People          People
	Houses          Housing
	Market          Market
	Companies       []Company
	Geography       Geography
}

func (s *Simulation) Tick() {
	nextDate := s.Date.AddDate(0, 0, int(s.SimulationSpeed))
	s.Date = nextDate
}

func (s *Simulation) ChangeSimulationSpeed() {
	switch s.SimulationSpeed {
	case Slow:
		s.SimulationSpeed = Mid
	case Mid:
		s.SimulationSpeed = Fast
	case Fast:
		s.SimulationSpeed = Pause
	default:
		s.SimulationSpeed = Slow
	}
}

func (s *Simulation) GetStats() string {
	return fmt.Sprintf("%s | Reserves: %s | Population: %4d (%+6.2f%%) | Free Houses: %2d | Unemployment: %5.2f%% | Companies: %2d | Market Value: %.2f (%+6.2f%%) | Inflation: %5.2f%%",
		s.Date.Format("2006-01-02"), utils.FormatCurrency(float64(s.Government.Reserves), "$"),
		s.People.Population, s.People.PopulationGrowthRate(), s.Houses.GetFreeHouses(),
		s.People.UnemploymentRate(), len(s.Companies), s.Market.GetMarketValue(),
		utils.GetLastValue(s.Market.History.MarketGrowthRate), utils.GetLastValue(s.Market.History.InflationRate))
}

var Sim Simulation

func NewSimulation(startYear int, houses int, governmentReserves int) Simulation {
	startDate := time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	return Simulation{
		SimulationSpeed: Mid,
		Date:            startDate,
		Government:      *NewGovernment(governmentReserves, startDate),
		People: People{
			Population:       0,
			PopulationValues: []int{0},
			LabourForce:      0,
			Unemployed:       0,
			Households:       []Household{},
		},
		Houses: *NewHousing(houses),
		Market: Market{
			InterestRate:           7.0,
			Unemployment:           0.001,
			LastCalculation:        startDate,
			MonthsOfNegativeGrowth: 0,
			History: MarketHistory{
				MarketValue:      []float64{1000},
				InflationRate:    []float64{0.001},
				MarketGrowthRate: []float64{0.001},
				MarketSentiment:  []float64{0.001},
				CompanyProfits:   []float64{0.001},
			},
		},
		Companies: []Company{},
		Geography: *NewGeography(64, 8, 3, 0.004, 0.01),
	}
}
