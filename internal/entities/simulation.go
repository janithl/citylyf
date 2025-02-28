package entities

import (
	"fmt"
	"maps"
	"slices"
	"sync"
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
	Mutex           sync.Mutex
	Government      Government
	People          People
	Houses          *Housing
	Market          Market
	Companies       map[int]*Company
	Geography       Geography
	tickNumber      int
}

func (s *Simulation) Tick() {
	s.tickNumber = (s.tickNumber + 100) % 1600
	if s.tickNumber%int(s.SimulationSpeed) == 0 {
		nextDate := s.Date.AddDate(0, 0, 1)
		s.Date = nextDate
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

func (s *Simulation) GetStats() string {
	return fmt.Sprintf("%s | Reserves: %s | Population: %d (%+06.2f%%) | Houses: %d (%d Free) | "+
		"Unemployment: %05.2f%% | Companies: %d | Market Value: %.2f (%+06.2f%%) | Inflation: %05.2f%% | IntRate: %05.2f%%",
		s.Date.Format("2006-01-02"), utils.FormatCurrency(float64(s.Government.Reserves), "$"), s.People.Population(),
		s.People.PopulationGrowthRate(), len(s.Houses.Houses), s.Houses.GetFreeHouses(), s.People.UnemploymentRate(),
		len(s.Companies), s.Market.MarketValue(), utils.GetLastValue(s.Market.History.MarketGrowthRate),
		s.Market.InflationRate(), s.Market.InterestRate())
}

// GetCompanyIDs returns a sorted list of company IDs
func (s *Simulation) GetCompanyIDs() []int {
	IDs := []int{}
	for company := range maps.Values(s.Companies) {
		IDs = append(IDs, company.ID)
	}
	slices.Sort(IDs)
	return IDs
}

func (s *Simulation) RegenerateMap(peakProb, rangeProb, cliffProb float64) {
	s.Geography = *NewGeography(64, 8, 3, peakProb, rangeProb, cliffProb)
}

var Sim Simulation

func NewSimulation(startYear, governmentReserves int) Simulation {
	startDate := time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	return Simulation{
		SimulationSpeed: Pause,
		Date:            startDate,
		Government:      *NewGovernment(governmentReserves, startDate),
		People: People{
			PopulationValues: []int{0},
			LabourForce:      0,
			Unemployed:       0,
			People:           make(map[int]*Person),
			Households:       make(map[int]*Household),
		},
		Houses: &Housing{
			LastHouseID: 100,
			Houses:      []House{},
		},
		Market: Market{
			NextRateRevision:       startDate.AddDate(0, 3, 0),
			MonthsOfNegativeGrowth: 0,
			History: MarketHistory{
				MarketValue:      []float64{1000},
				InflationRate:    []float64{0.001},
				InterestRate:     []float64{7.0},
				MarketGrowthRate: []float64{0.001},
				MarketSentiment:  []float64{0.001},
				CompanyProfits:   []float64{0.001},
			},
		},
		Companies: make(map[int]*Company),
		Geography: *NewGeography(64, 8, 3, 0.0015, 0.005, 0.01),
	}
}
