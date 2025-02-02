package economy

import (
	"math"
	"math/rand"
)

type Market struct {
	InterestRate       float64 // Higher interest rates = bad for stocks
	LastInflationRate  float64 // High inflation leads to money tightening
	Unemployment       float64 // High unemployment reduces market confidence
	CorporateTax       float64 // High corporate tax = lower stock growth
	GovernmentSpending float64 // More spending = Higher money supply
}

// Random factor (news, global events)
func (m *Market) MarketSentiment() float64 {
	return (rand.Float64() * 4) - 2 // Random factor (-2% to +2%)
}

// Supply-chain disruptions (0 = normal, higher = worse)
func (m *Market) SupplyShock() float64 {
	return rand.Float64() * 3 // 0% - 3%
}

const baseGrowth = 3.0 // Normal money supply growth (%)

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
	totalGrowth := baseGrowth + interestImpact + inflationImpact + spendingImpact + confidenceImpact

	// Ensure realistic money supply growth
	if totalGrowth < 0 {
		totalGrowth = 0 // Contraction capped at 0%
	}
	if totalGrowth > 15 {
		totalGrowth = 15 // Hyper-expansion capped at 15%
	}

	return totalGrowth
}

const baseInflation = 2.0 // Normal baseline inflation (%)

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
	totalInflation := baseInflation + moneyImpact + interestImpact + demandImpact + supplyImpact

	// Ensure inflation stays reasonable
	if totalInflation < -1 {
		totalInflation = -1 // Deflation capped at -1%
	}
	if totalInflation > 20 {
		totalInflation = 20 // Hyperinflation capped at 20%
	}

	return totalInflation
}
