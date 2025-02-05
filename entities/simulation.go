package entities

import (
	"fmt"
	"time"
)

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

func (s *Simulation) ChangeSimulationSpeed() {
	switch s.SimulationSpeed {
	case Slow:
		s.SimulationSpeed = Mid
	case Mid:
		s.SimulationSpeed = Fast
	default:
		s.SimulationSpeed = Slow
	}
}

func (s *Simulation) GetStats() string {
	return fmt.Sprintf("%s | Population: %4d (%+6.2f%%) | Unemployment: %5.2f%% | Companies: %d | Market Value: %.2f (%+6.2f%%) | Inflation: %5.2f%%",
		s.Date.Format("2006-01-02"), s.People.Population, s.People.PopulationGrowthRate(), s.People.UnemploymentRate(),
		len(s.Companies), s.Market.GetMarketValue(), s.Market.LastMarketGrowthRate, s.Market.LastInflationRate)
}

var Sim Simulation

func NewSimulation(startYear int) Simulation {
	startDate := time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	return Simulation{
		SimulationSpeed: Mid,
		Date:            startDate,
		People: People{
			Population:       0,
			PopulationValues: []int{0},
			LabourForce:      0,
			Unemployed:       0,
			Households:       []Household{},
		},
		Market: Market{
			InterestRate:       7.0,
			Unemployment:       0.001,
			CorporateTax:       2.0,
			GovernmentSpending: 5.0,

			LastCalculation:        startDate,
			LastInflationRate:      0.001,
			LastMarketGrowthRate:   0.001,
			LastMarketSentiment:    0.001,
			LastSixMonthsProfits:   []float64{0.0},
			MarketHigh:             1000,
			MarketValues:           []float64{1000},
			MonthsOfNegativeGrowth: 0,
		},
		Companies: []Company{},
	}
}
