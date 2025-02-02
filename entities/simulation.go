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
}

func (s *Simulation) Tick() {
	nextDate := s.Date.AddDate(0, 0, int(s.SimulationSpeed))
	s.Date = nextDate
}

func NewSimulation(startYear int) Simulation {
	return Simulation{
		SimulationSpeed: Mid,
		Date:            time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}
