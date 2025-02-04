package entities

type People struct {
	Population       int
	PopulationValues []int // Historical values
	PopulationHigh   int   // Historic maximum poplulation
	LabourForce      int   // Employable people
	Unemployed       int
	Households       []Household
}

func (p *People) UnemploymentRate() float64 {
	return 100.0 * float64(p.Unemployed) / float64(p.LabourForce)
}

func (p *People) PopulationGrowthRate() float64 {
	lastPopulationValue := p.PopulationValues[len(p.PopulationValues)-1]
	return 100.0 * float64(p.Population-lastPopulationValue) / float64(lastPopulationValue)
}

func (p *People) MoveIn(h Household) {
	p.Households = append(p.Households, h)
	p.Population += len(h.Members)
}

// calculate the unemployed and the total labour force
func (p *People) CalculateUnemployment() {
	labourforce, unemployed := 0, 0
	for i := 0; i < len(p.Households); i++ {
		for j := 0; j < len(p.Households[i].Members); j++ {
			if p.Households[i].Members[j].IsEmployable() {
				labourforce += 1
				if !p.Households[i].Members[j].IsEmployed() {
					unemployed += 1
				}
			}
		}
	}
	p.LabourForce = labourforce
	p.Unemployed = unemployed
}

// Append current population value to history
func (p *People) UpdatePopulationValues() {
	if len(p.PopulationValues) >= 10 {
		p.PopulationValues = p.PopulationValues[1:] // Remove first element (FIFO behavior)
	}
	p.PopulationValues = append(p.PopulationValues, p.Population)
	if p.Population > p.PopulationHigh { // Set population high
		p.PopulationHigh = p.Population
	}
}
