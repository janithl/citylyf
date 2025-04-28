package entities

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

// SimulationSpeed defines the speed of the simulation
type SimulationSpeed int

const (
	Pause SimulationSpeed = 0
	Slow  SimulationSpeed = 1600
	Mid   SimulationSpeed = 400
	Fast  SimulationSpeed = 100
)

type Simulation struct {
	SimulationSpeed SimulationSpeed
	Date            time.Time
	Mutex           sync.RWMutex
	Government      *Government
	People          *People
	Houses          Housing
	Companies       Companies
	Market          *Market
	Geography       *Geography
	tickNumber      int
	lastID          atomic.Uint32
	SavePath        string
	NameService     *NameService
}

func (s *Simulation) Tick(dailyActivity func()) {
	s.tickNumber = (s.tickNumber + 100) % 1600
	if s.tickNumber%int(s.SimulationSpeed) == 0 {
		nextDate := s.Date.AddDate(0, 0, 1)
		s.Date = nextDate
		dailyActivity()
	}
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

func (s *Simulation) PauseSimulation() {
	s.SimulationSpeed = Pause
}

func (s *Simulation) GetStats() string {
	return fmt.Sprintf("%s | Reserves: %s | Population: %d (%+06.2f%%) | Houses: %d (%d Free) | "+
		"Unemployment: %05.2f%% | Companies: %d | Market Value: %.2f (%+06.2f%%) | Inflation: %05.2f%% | IntRate: %05.2f%%",
		s.Date.Format("2006-01-02"), utils.FormatCurrency(s.Government.GetReservesAtHand(), "$"), s.People.Population(),
		s.People.PopulationGrowthRate(), len(s.Houses), s.Houses.GetFreeHouses(), s.People.UnemploymentRate(),
		len(s.Companies), s.Market.MarketValue(), utils.GetLastValue(s.Market.History.MarketGrowthRate),
		s.Market.InflationRate(), s.Market.InterestRate())
}

func (s *Simulation) GetNextID() int {
	return int(s.lastID.Add(1))
}

func (s *Simulation) SendStats() {
	select {
	case SimStats <- s.GetStats():
	default:
	}
}

func (s *Simulation) RegenerateMap(peakProb, rangeProb, cliffProb float64) {
	s.Geography = NewGeography(64, 8, 8, 3, 7, peakProb, rangeProb, cliffProb)
}

var Sim *Simulation
var SimStats chan string

func NewSimulation(startYear, governmentReserves int) *Simulation {
	startDate := time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	sim := &Simulation{
		SimulationSpeed: Pause,
		Date:            startDate,
		Government:      NewGovernment(governmentReserves, startDate),
		People: &People{
			LabourForce:            0,
			Unemployed:             0,
			PopulationValues:       []int{0},
			UnemploymentRateValues: []float64{0.0},
			AverageWageValues:      []float64{0.0},
			People:                 make(map[int]*Person),
			Households:             make(map[int]*Household),
		},
		Houses:    make(map[int]*House),
		Companies: make(map[int]*Company),
		Market: &Market{
			NextRateRevision:       startDate.AddDate(0, 3, 0),
			MonthsOfNegativeGrowth: 0,
			History: MarketHistory{
				MarketValue:      []float64{1000},
				InflationRate:    []float64{0.001},
				InterestRate:     []float64{7.0},
				MarketGrowthRate: []float64{0.001},
				MarketSentiment:  []float64{0.001},
				CompanyProfits:   []float64{0.001},
				AverageRent:      []float64{0.0},
			},
		},
		Geography:   NewGeography(64, 8, 8, 3, 7, 0.0015, 0.005, 0.01),
		NameService: NewNameService(),
	}
	sim.lastID.Store(10000)         // start IDs at 10000
	SimStats = make(chan string, 1) // create the stats channel

	return sim
}

func LoadSimulationFromSave(path string, sim *Simulation, lastID uint32, tiles [][]Tile, roads []*Road) {
	Sim = sim
	Sim.lastID.Store(lastID)
	Sim.SavePath = path

	Sim.Geography.tiles = tiles
	Sim.Geography.roads = roads
	SimStats = make(chan string, 1)
}

func StartNewSim() {
	Sim = NewSimulation(2020, 1000000)
	Sim.SendStats()
}
