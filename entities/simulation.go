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
	CityData        CityData
}

func (s *Simulation) Tick() {
	nextDate := s.Date.AddDate(0, 0, int(s.SimulationSpeed))
	s.Date = nextDate
}

var CitySimulation Simulation

func NewSimulation(startYear int) Simulation {
	return Simulation{
		SimulationSpeed: Mid,
		Date:            time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC),
		CityData: CityData{
			Population: 0,
			Households: []Household{},
		},
	}
}
