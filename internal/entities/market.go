package entities

import (
	"math"
	"math/rand"
	"time"

	"github.com/janithl/citylyf/internal/utils"
)

const ( // Constants for economic behavior
	BaseMoneySupplyGrowth = 3.0
	BaseInflation         = 2.0
	BaseMarketGrowth      = 2.0
	MinGrowth             = -5.0
	RecessionThreshold    = -2.0 // Growth below this triggers a recession
	BoomThreshold         = 5.0  // Growth above this triggers a boom
)

type MarketHistory struct { // tracking last 12 months data
	MarketValue, InflationRate, MarketGrowthRate, MarketSentiment, CompanyProfits []float64
}

// Market tracks economic cycles and financial conditions
type Market struct {
	InterestRate, Unemployment float64
	LastCalculation            time.Time
	History                    MarketHistory
	MonthsOfNegativeGrowth     int
	InRecession, InBoom        bool
}

// MarketSentiment adjusts sentiment based on boom/bust cycles
func (m *Market) MarketSentiment() float64 {
	baseSentiment := (rand.Float64() * 4) - 2 // Random factor (-2% to +2%)

	if m.InRecession { // Modify Sentiment Based on Boom/Bust Cycle
		baseSentiment -= (rand.Float64() * 2) // Negative bias (-0% to -2% extra)
	} else if m.InBoom {
		baseSentiment += (rand.Float64() * 2) // Positive bias (+0% to +2% extra)
	}

	if baseSentiment < -3 { // Clamp sentiment to a reasonable range**
		baseSentiment = -3
	} else if baseSentiment > 3 {
		baseSentiment = 3
	}

	m.History.MarketSentiment = utils.AddFifo(m.History.MarketSentiment, baseSentiment, 10)
	return baseSentiment
}

// SupplyShock applies supply-chain disruptions (0% - 3%)
func (m *Market) SupplyShock() float64 {
	return rand.Float64() * 3
}

// MoneySupplyGrowth calculates money supply changes
func (m *Market) MoneySupplyGrowth() float64 {
	lastInflationRate := utils.GetLastValue(m.History.InflationRate)
	interestImpact := -math.Pow(m.InterestRate/5, 1.2)             // High rates slow money supply
	inflationImpact := -math.Pow((lastInflationRate-2)/4, 2)       // High inflation slows supply
	spendingImpact := Sim.Government.GetGovernmentSpending() * 0.5 // More spending increases supply
	confidenceImpact := m.MarketSentiment() * 0.3                  // Market sentiment effect
	totalGrowth := BaseMoneySupplyGrowth + interestImpact + inflationImpact + spendingImpact + confidenceImpact

	if totalGrowth < 0 {
		totalGrowth = 0 // Cap contraction at 0%
	} else if totalGrowth > 10 {
		totalGrowth = 10 // Cap expansion at 15%
	}
	return totalGrowth
}

// Inflation calculates inflation considering money supply, interest rates, and supply shocks
func (m *Market) Inflation(populationGrowth float64) float64 {
	moneyImpact := math.Log(m.MoneySupplyGrowth()+1) * 1.5 // More money = higher inflation
	interestImpact := -math.Pow(m.InterestRate/3, 1.5)     // Higher rates reduce inflation
	demandImpact := math.Max(populationGrowth*0.5, 0.0)    // Higher demand pushes inflation up
	supplyImpact := m.SupplyShock() * 1.2                  // Supply disruptions worsen inflation
	totalInflation := BaseInflation + moneyImpact + interestImpact + demandImpact + supplyImpact

	if totalInflation < -1 {
		totalInflation = -1 // Cap deflation at -1%
	} else if totalInflation > 15 {
		totalInflation = 15 // Cap hyperinflation at 15%
	}

	m.History.InflationRate = utils.AddFifo(m.History.InflationRate, totalInflation, 10)
	return totalInflation
}

// MarketGrowth calculates stock index growth with boom/bust cycle logic
func (m *Market) MarketGrowth() float64 {
	lastInflationRate := utils.GetLastValue(m.History.InflationRate)
	lastMarketGrowthRate := utils.GetLastValue(m.History.MarketGrowthRate)

	interestImpact := -math.Pow(m.InterestRate/8, 2)                       // Higher rates slow growth
	inflationImpact := -math.Pow((lastInflationRate-5)/3, 2)               // Inflation impact (good at 3-5%, bad above 6%)
	unemploymentImpact := -m.Unemployment / 25                             // High unemployment reduces spending
	taxImpact := -(Sim.Government.CorporateTaxRate) / 30                   // Higher taxes = lower market growth
	marketSentimentImpact := utils.GetLastValue(m.History.MarketSentiment) // External random factors
	profitImpact := m.calculateProfitImpact()                              // Effect of corporate profits

	recoveryBoost := 0.0
	if m.MonthsOfNegativeGrowth > 3 { // If negative for 3+ months, slow recovery starts
		recoveryBoost = float64(m.MonthsOfNegativeGrowth) * 0.5
	}

	// **Boom/Bust Cycle Mechanic**
	if m.InRecession && lastMarketGrowthRate > RecessionThreshold { // Exit recession if market improves
		m.InRecession = false
	}
	if m.InBoom && lastMarketGrowthRate < BoomThreshold { // Exit boom if market slows
		m.InBoom = false
	}
	if lastMarketGrowthRate < RecessionThreshold && !m.InRecession { // Enter recession if growth is too low
		m.InRecession = true
	}
	if lastMarketGrowthRate > BoomThreshold && !m.InBoom { // Enter boom if growth is too high
		m.InBoom = true
	}

	cycleImpact := 0.0
	if m.InRecession {
		cycleImpact = -1.0 + (rand.Float64() * 1.5) // Mild drag with jitter
	}
	if m.InBoom {
		cycleImpact = 2.5 + (rand.Float64() * 1.0) // Strong boost with jitter
	}

	longTermCorrection := (BaseMarketGrowth - lastMarketGrowthRate) * 0.1 // correction to avoid market collapse

	totalGrowth := BaseMarketGrowth + interestImpact + inflationImpact + unemploymentImpact + taxImpact +
		marketSentimentImpact + profitImpact + recoveryBoost + cycleImpact + longTermCorrection

	// Prevent extreme crashes
	if totalGrowth < MinGrowth {
		totalGrowth = MinGrowth
	}

	m.History.MarketGrowthRate = utils.AddFifo(m.History.MarketGrowthRate, totalGrowth, 10)
	return totalGrowth
}

// calculateProfitImpact smooths corporate profit influence on the market
func (m *Market) calculateProfitImpact() float64 {
	totalProfit := 0.0
	for _, p := range m.History.CompanyProfits {
		totalProfit += p
	}
	averageProfit := totalProfit / float64(len(m.History.CompanyProfits))

	if averageProfit > 0 {
		return math.Log1p(averageProfit) / 10 // Moderate growth effect
	}
	return -math.Sqrt(math.Abs(averageProfit)) / 50 // Controlled negative impact
}

// GetMarketValue returns the latest market index value
func (m *Market) GetMarketValue() float64 {
	return utils.GetLastValue(m.History.MarketValue)
}

// UpdateMarketValue updates market history & records highs
func (m *Market) UpdateMarketValue(marketGrowth float64) float64 {
	lastMarketValue := m.GetMarketValue()
	newMarketValue := lastMarketValue + (lastMarketValue * marketGrowth / 100)

	m.History.MarketValue = utils.AddFifo(m.History.MarketValue, newMarketValue, 10)
	m.LastCalculation = Sim.Date

	return newMarketValue
}

// Total company profits are added to market history
func (m *Market) ReportCompanyProfits(profits float64) {
	m.History.CompanyProfits = utils.AddFifo(m.History.CompanyProfits, profits, 10)
}
