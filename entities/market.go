package entities

import (
	"math"
	"math/rand"
	"time"
)

const BaseMoneySupplyGrowth = 3.0 // Normal money supply growth (%)
const BaseInflation = 2.0         // Normal baseline inflation (%)
const BaseMarketGrowth = 2.0      // Expected annual growth in % (adjusted by economic conditions)
const MinGrowth = -10.0           // Minimum allowed negative growth

type Market struct {
	InterestRate       float64 // Higher interest rates = bad for stocks
	Unemployment       float64 // High unemployment reduces market confidence
	CorporateTax       float64 // High corporate tax = lower stock growth
	GovernmentSpending float64 // More spending = Higher money supply

	// Historical
	LastCalculation        time.Time // last time the market was calculated
	LastInflationRate      float64   // High inflation leads to money tightening
	LastMarketGrowthRate   float64   // Growth adjusts the stock index based on economic conditions
	LastMarketSentiment    float64   // Random factor (news, global events)
	LastSixMonthsProfits   []float64
	MarketHigh             float64   // The highest value the market has recorded
	MarketValues           []float64 // The historical market values
	MonthsOfNegativeGrowth int       // Holds number of months of negative growth for recession
}

// Random factor (news, global events)
func (m *Market) MarketSentiment() float64 {
	sentiment := (rand.Float64() * 4) - 2 // Random factor (-2% to +2%)
	m.LastMarketSentiment = sentiment
	return sentiment
}

// Supply-chain disruptions (0 = normal, higher = worse)
func (m *Market) SupplyShock() float64 {
	return rand.Float64() * 3 // 0% - 3%
}

// MoneySupplyGrowth calculates money supply growth based on economic conditions
func (m *Market) MoneySupplyGrowth() float64 {

	// Interest Rate Effect: Higher rates slow money growth
	interestImpact := -math.Pow(m.InterestRate/5, 1.2)

	// Inflation Effect: High inflation (>6%) slows money supply to prevent overheating
	inflationImpact := -math.Pow((m.LastInflationRate-2)/4, 2)

	// Government Spending Effect: More spending increases money supply
	spendingImpact := m.GovernmentSpending * 0.5

	// Market Sentiment Effect: High confidence leads to more liquidity
	confidenceImpact := m.MarketSentiment() * 0.3

	// Calculate total money supply growth
	totalGrowth := BaseMoneySupplyGrowth + interestImpact + inflationImpact + spendingImpact + confidenceImpact

	// Ensure realistic money supply growth
	if totalGrowth < 0 {
		totalGrowth = 0 // Contraction capped at 0%
	}
	if totalGrowth > 15 {
		totalGrowth = 15 // Hyper-expansion capped at 15%
	}

	return totalGrowth
}

// Inflation calculates inflation based on economic conditions
func (m *Market) Inflation(populationGrowth float64) float64 {
	// Money supply impact: More money = more inflation (logarithmic effect)
	moneyImpact := math.Log(m.MoneySupplyGrowth()+1) * 1.5

	// Interest Rate Effect: Higher rates slow inflation
	interestImpact := -math.Pow(m.InterestRate/5, 1.2)

	// Demand Growth: Higher demand pushes inflation up
	demandImpact := math.Max(populationGrowth*0.8, 0.0)
	if populationGrowth > 5 {
		demandImpact = populationGrowth * 0.8
	}

	// Supply Shock: Disruptions push inflation higher
	supplyImpact := m.SupplyShock() * 1.2

	// Total inflation calculation
	totalInflation := BaseInflation + moneyImpact + interestImpact + demandImpact + supplyImpact

	// Ensure inflation stays reasonable
	if totalInflation < -1 {
		totalInflation = -1 // Deflation capped at -1%
	}
	if totalInflation > 20 {
		totalInflation = 20 // Hyperinflation capped at 20%
	}

	m.LastInflationRate = totalInflation
	return totalInflation
}

// MarketGrowth adjusts the stock index based on economic conditions
func (m *Market) MarketGrowth() float64 {
	// Interest Rate Effect: High rates slow down growth
	interestImpact := -math.Pow(m.InterestRate/5, 2)

	// Inflation Effect: Moderate (2-3%) is good, but high (>6%) is bad
	inflationImpact := -math.Pow((m.LastInflationRate-2)/3, 2)

	// Unemployment Effect: More unemployment → Less consumer spending → Weaker market
	unemploymentImpact := -m.Unemployment / 10

	// Corporate Tax Effect: Higher tax rates = Lower stock growth
	taxImpact := -m.CorporateTax / 30

	// Market Sentiment: Random external factors (news, global events, speculation)
	marketSentimentImpact := m.LastMarketSentiment

	// Instead of using last month's profit, calculate the rolling average over 6 months.
	totalProfit := 0.0
	for _, p := range m.LastSixMonthsProfits {
		totalProfit += p
	}
	averageProfit := totalProfit / float64(len(m.LastSixMonthsProfits))

	var profitImpact float64
	if averageProfit > 0 {
		profitImpact = math.Log1p(averageProfit) / 10 // Normalized positive profit effect
	} else {
		profitImpact = -math.Sqrt(math.Abs(averageProfit)) / 50 // Controlled negative impact
	}

	// Market Recovery Mechanism
	recoveryBoost := 0.0
	if m.MonthsOfNegativeGrowth > 3 { // If negative for 3+ months, slow recovery kicks in
		recoveryBoost = float64(m.MonthsOfNegativeGrowth) * 0.2
	}

	// Calculate total stock market growth
	totalGrowth := BaseMoneySupplyGrowth + interestImpact + inflationImpact + unemploymentImpact + taxImpact + marketSentimentImpact + profitImpact + recoveryBoost

	// Update rolling history of profits
	m.LastSixMonthsProfits = append(m.LastSixMonthsProfits, averageProfit)
	if len(m.LastSixMonthsProfits) > 6 {
		m.LastSixMonthsProfits = m.LastSixMonthsProfits[1:] // Keep last 6 months only
	}

	// Reset negative growth counter if market turns positive
	if totalGrowth > 0 {
		m.MonthsOfNegativeGrowth = 0
	} else {
		m.MonthsOfNegativeGrowth++
	}

	// Prevent extreme crashes (below MinGrowth)
	if totalGrowth < MinGrowth {
		totalGrowth = MinGrowth
	}

	m.LastMarketGrowthRate = totalGrowth
	return totalGrowth
}

func (m *Market) GetMarketValue() float64 {
	return m.MarketValues[len(m.MarketValues)-1]
}

// Append new market value to history
func (m *Market) UpdateMarketValue(marketGrowth float64) float64 {
	if len(m.MarketValues) >= 10 {
		m.MarketValues = m.MarketValues[1:] // Remove first element (FIFO behavior)
	}
	lastMarketValue := m.MarketValues[len(m.MarketValues)-1]
	newMarketValue := lastMarketValue + (lastMarketValue * marketGrowth / 100)
	if newMarketValue > m.MarketHigh { // Set market high
		m.MarketHigh = newMarketValue
	}

	m.MarketValues = append(m.MarketValues, newMarketValue) // append new market value
	m.LastCalculation = Sim.Date                            // update last calculation time

	return newMarketValue
}
