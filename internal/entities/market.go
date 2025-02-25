package entities

import (
	"fmt"
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
	InflationTarget       = 2.0
)

type MarketHistory struct { // tracking last 12 months data
	MarketValue, InflationRate, InterestRate, MarketGrowthRate, MarketSentiment, CompanyProfits []float64
}

// Market tracks economic cycles and financial conditions
type Market struct {
	NextRateRevision       time.Time
	History                MarketHistory
	MonthsOfNegativeGrowth int
	InRecession, InBoom    bool
}

func (m *Market) InterestRate() float64 {
	return utils.GetLastValue(m.History.InterestRate)
}

func (m *Market) InflationRate() float64 {
	return utils.GetLastValue(m.History.InflationRate)
}

func (m *Market) MarketValue() float64 {
	return utils.GetLastValue(m.History.MarketValue)
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
	interestImpact := -math.Pow(m.InterestRate()/5, 1.2)           // High rates slow money supply
	inflationImpact := -math.Pow((m.InflationRate()-2)/4, 2)       // High inflation slows supply
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

// CalculateInflation calculates inflation considering money supply, interest rates, and supply shocks
func (m *Market) CalculateInflation(populationGrowth float64) {
	moneyImpact := math.Log(m.MoneySupplyGrowth()+1) * 1.5 // More money = higher inflation
	interestImpact := -math.Pow(m.InterestRate()/3, 1.5)   // Higher rates reduce inflation
	demandImpact := math.Max(populationGrowth*0.5, 0.0)    // Higher demand pushes inflation up
	supplyImpact := m.SupplyShock() * 1.2                  // Supply disruptions worsen inflation
	totalInflation := BaseInflation + moneyImpact + interestImpact + demandImpact + supplyImpact

	if totalInflation < -1 {
		totalInflation = -1 // Cap deflation at -1%
	} else if totalInflation > 15 {
		totalInflation = 15 // Cap hyperinflation at 15%
	}

	m.History.InflationRate = utils.AddFifo(m.History.InflationRate, totalInflation, 10)
}

// CalculateMarketGrowth calculates stock index growth with boom/bust cycle logic
func (m *Market) CalculateMarketGrowth() float64 {
	lastMarketGrowthRate := utils.GetLastValue(m.History.MarketGrowthRate)

	interestImpact := -math.Pow(m.InterestRate()/8, 2)                     // Higher rates slow growth
	inflationImpact := -math.Pow((m.InflationRate()-5)/3, 2)               // Inflation impact (good at 3-5%, bad above 6%)
	unemploymentImpact := -Sim.People.UnemploymentRate() / 25              // High unemployment reduces spending
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
	if len(m.History.CompanyProfits) < 1 {
		return 0.0
	}

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

// UpdateMarketValue updates market history & records highs
func (m *Market) UpdateMarketValue(marketGrowth float64) float64 {
	lastMarketValue := m.MarketValue()
	newMarketValue := lastMarketValue + (lastMarketValue * marketGrowth / 100)
	m.History.MarketValue = utils.AddFifo(m.History.MarketValue, newMarketValue, 10)
	return newMarketValue
}

// Total company profits are added to market history
func (m *Market) ReportCompanyProfits(profits float64) {
	m.History.CompanyProfits = utils.AddFifo(m.History.CompanyProfits, profits, 10)
}

// ReviseInterestRate updates interest rate based on inflation
func (m *Market) ReviseInterestRate() {
	if Sim.Date.Before(m.NextRateRevision) || len(m.History.InflationRate) < 3 {
		return // rate revisions only happen once a quarter + we need historical inflation data to do a rates revision
	}
	averageInflationRate := 0.0
	for i := len(m.History.InflationRate) - 3; i < len(m.History.InflationRate); i++ {
		averageInflationRate += m.History.InflationRate[i]
	}
	averageInflationRate /= 3

	interestRateChange := 0.0
	switch {
	case averageInflationRate < InflationTarget-1.5:
		// large deflation is happening, lower interest rates by 50 basis points
		interestRateChange = -0.5
	case averageInflationRate < InflationTarget-0.75:
		// small deflation is happening, lower interest rates by 25 basis points
		interestRateChange = -0.25
	case averageInflationRate > InflationTarget+0.75:
		// small inflation is happening, raise interest rates by 25 basis points
		interestRateChange = 0.25
	case averageInflationRate > InflationTarget+1.5:
		// large inflation is happening, raise interest rates by 50 basis points
		interestRateChange = 0.5
	}

	newInterestRate := m.InterestRate() + interestRateChange
	if newInterestRate < 0 { // we cannot allow negative interest rates
		return
	}
	m.History.InterestRate = utils.AddFifo(m.History.InterestRate, newInterestRate, 20)
	m.NextRateRevision = Sim.Date.AddDate(0, 3, 0) // next rate revision in 3 months

	if interestRateChange > 0 {
		fmt.Printf("[ Rate ] Avg. inflation at %.2f%%, above target range. Interest rate raised by %.2f%% to", averageInflationRate, interestRateChange)
	} else if interestRateChange < 0 {
		fmt.Printf("[ Rate ] Avg. inflation at %.2f%%, below target range. Interest rate lowered by %.2f%% to", averageInflationRate, interestRateChange)
	} else {
		fmt.Printf("[ Rate ] Avg. inflation at %.2f%%, within the target range. Interest rates held steady at", averageInflationRate)
	}
	fmt.Printf(" %.2f%%. Next rates revision on %s\n", newInterestRate, m.NextRateRevision.Format("2006-01-02"))
}
