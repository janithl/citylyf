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
	SevereUnemployment    = 10.0 // 10% unemployment is considered severe
)

type MarketHistory struct { // tracking last 12 months data
	MarketValue, InflationRate, InterestRate, MarketGrowthRate, MarketSentiment, CompanyProfits, AverageRent []float64
}

// Market tracks economic cycles and financial conditions
type Market struct {
	NextRateRevision            time.Time
	History                     MarketHistory
	MonthsOfNegativeGrowth      int
	InRecession, InBoom         bool
	HousingDemand, RetailDemand float64
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

	baseSentiment = utils.Clamp(baseSentiment, -3, 3) // Clamp sentiment to a reasonable range**
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

	return utils.Clamp(totalGrowth, 0, 15) // Cap contraction at 0% and  expansion at 15%
}

// CalculateInflation calculates inflation considering money supply, interest rates, and supply shocks
func (m *Market) CalculateInflation(populationGrowth float64) {
	moneyImpact := math.Log(m.MoneySupplyGrowth()+1) * 1.5 // More money = higher inflation
	interestImpact := -math.Pow(m.InterestRate()/3, 1.5)   // Higher rates reduce inflation
	demandImpact := math.Max(populationGrowth*0.5, 0.0)    // Higher demand pushes inflation up
	supplyImpact := m.SupplyShock() * 1.2                  // Supply disruptions worsen inflation
	totalInflation := BaseInflation + moneyImpact + interestImpact + demandImpact + supplyImpact

	totalInflation = utils.Clamp(totalInflation, -1, 15) // Cap deflation at -1% and hyperinflation at 15%
	m.History.InflationRate = utils.AddFifo(m.History.InflationRate, totalInflation, 20)
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
	m.History.MarketValue = utils.AddFifo(m.History.MarketValue, newMarketValue, 20)
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

	// set rate change, working within inflation target band
	var interestRateChange float64
	if averageInflationRate < InflationTarget-1.0 { // inflation lower than target band, lower rates
		interestRateChange = (averageInflationRate - InflationTarget) * 0.5
	} else if averageInflationRate > InflationTarget+1.0 { // inflation is too high, raise rates proportionally
		interestRateChange = (averageInflationRate - InflationTarget) * 0.25
	} else { // Inflation is within the target band
		// Only lower rates if the previous revision was held steady.
		if len(m.History.InterestRate) > 1 {
			lastRate := m.History.InterestRate[len(m.History.InterestRate)-1]
			secondLastRate := m.History.InterestRate[len(m.History.InterestRate)-2]
			// Use a small epsilon for float comparison.
			if math.Abs(lastRate-secondLastRate) < 1e-6 {
				interestRateChange = -0.25
			} else {
				interestRateChange = 0.0
			}
		} else { // hold rates steady
			interestRateChange = 0.0
		}
	}

	interestRateChange = math.Round(interestRateChange*4) / 4       // always change rates in multiples of 0.25%
	interestRateChange = utils.Clamp(interestRateChange, -0.5, 0.5) // Cap the change to avoid overshooting (max +/- 0.5% per revision).

	// Calculate the new interest rate.
	newInterestRate := m.InterestRate() + interestRateChange
	if newInterestRate < 0 { // we cannot allow negative interest rates
		newInterestRate = 0
	}

	m.History.InterestRate = utils.AddFifo(m.History.InterestRate, newInterestRate, 20)
	m.NextRateRevision = Sim.Date.AddDate(0, 3, 0) // next rate revision in 3 months

	if interestRateChange > 0 {
		fmt.Printf("[ Rate ] Avg. inflation at %.2f%%, above target range. Interest rate raised by %.2f%% to", averageInflationRate, interestRateChange)
	} else if interestRateChange < 0 {
		fmt.Printf("[ Rate ] Avg. inflation at %.2f%%, within or below target range. Interest rate lowered by %.2f%% to", averageInflationRate, interestRateChange)
	} else {
		fmt.Printf("[ Rate ] Avg. inflation at %.2f%%, within the target range. Interest rates held steady at", averageInflationRate)
	}
	fmt.Printf(" %.2f%%. Next rates revision on %s\n", newInterestRate, m.NextRateRevision.Format("2006-01-02"))
}

func (m *Market) CalculateHousingAndRetailDemand(totalHouses, vacantHouses int) {
	housingDemand, retailDemand := 0.5, 0.5 // Neutral starting point

	// **Housing Demand Factors**
	populationGrowth := Sim.People.PopulationGrowthRate()            // Higher population growth increases housing demand
	unemploymentRate := Sim.People.UnemploymentRate()                // Higher unemployment reduces ability to afford housing
	incomeGrowth := Sim.Market.CalculateMarketGrowth() * 0.2         // Higher market growth generally leads to better wages
	interestRate := Sim.Market.InterestRate()                        // Higher rates make mortgages more expensive
	vacancyImpact := float64(vacantHouses) / float64(totalHouses+1)  // Housing Availability Impact, +1 prevents div by zero
	minDemandBoost := 1.0 / (1.0 + float64(Sim.People.Population())) // Creates demand when population is near 0

	housingDemand += (populationGrowth * 2)  // Direct impact of growth
	housingDemand -= (unemploymentRate / 20) // Unemployment suppresses demand
	housingDemand += (incomeGrowth / 10)     // More income supports higher demand
	housingDemand -= (interestRate / 10)     // High interest rates reduce borrowing
	housingDemand -= vacancyImpact           // More vacant houses reduce demand
	housingDemand += minDemandBoost          // Ensures some demand at low population

	// **Retail Demand Factors**
	disposableIncome := Sim.People.AverageMonthlyDisposableIncome()     // Higher disposable income increases retail demand
	consumerConfidence := utils.GetLastValue(m.History.MarketSentiment) // Market sentiment affects spending habits
	taxImpact := -Sim.Government.SalesTaxRate / 20                      // Higher sales tax slightly reduces demand
	unemploymentImpact := -(unemploymentRate / 15)                      // Unemployment reduces disposable income
	profitImpact := m.calculateProfitImpact()                           // Stronger corporate profits usually reflect strong consumer demand

	retailDemand += float64(disposableIncome / 50000) // Normalize disposable income impact
	retailDemand += (consumerConfidence / 10)         // Positive market sentiment encourages spending
	retailDemand += profitImpact                      // Profitable businesses suggest strong demand
	retailDemand -= taxImpact                         // Higher taxes suppress demand slightly
	retailDemand -= unemploymentImpact                // Unemployment hurts demand

	// Clamp values between 0 and 1
	m.HousingDemand = utils.Clamp(housingDemand, 0, 1)
	m.RetailDemand = utils.Clamp(retailDemand, 0, 1)
	fmt.Printf("[ Econ ] Housing demand is at %.2f and Retail demand is at %.2f\n", m.HousingDemand, m.RetailDemand)
}
