package entities

import "time"

// SimulationSpeed defines the number of days at which the sim moves
type SimulationSpeed int

const (
	Slow SimulationSpeed = 3
	Mid  SimulationSpeed = 7
	Fast SimulationSpeed = 28
)

type Simulation struct {
	SimulationSpeed SimulationSpeed
	Date            time.Time
	People          People
	Market          Market
	Companies       []Company
}

func (s *Simulation) Tick() {
	nextDate := s.Date.AddDate(0, 0, int(s.SimulationSpeed))
	s.Date = nextDate
}

var Sim Simulation

func NewSimulation(startYear int) Simulation {
	startDate := time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	return Simulation{
		SimulationSpeed: Mid,
		Date:            startDate,
		People: People{
			Population: 0,
			Households: []Household{},
		},
		Market: Market{
			InterestRate:           5.0,
			LastInflationRate:      0.0,
			Unemployment:           0.0,
			CorporateTax:           3.0,
			GovernmentSpending:     10.0,
			MonthsOfNegativeGrowth: 0,
			LastCalculation:        startDate,
			MarketHigh:             1000,
			MarketValues:           []float64{1000},
		},
		Companies: []Company{},
	}
}
