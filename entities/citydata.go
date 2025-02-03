package entities

type People struct {
	Population  int
	LabourForce int
	Unemployed  int
	Households  []Household
}

func (c *People) UnemploymentRate() float64 {
	return 100.0 * float64(c.Unemployed) / float64(c.LabourForce)
}

func (c *People) MoveIn(h Household) {
	c.Households = append(c.Households, h)
	c.Population += len(h.Members)
}

func (c *People) CalculateUnemployment() {
	labourforce, unemployed := 0, 0
	for i := 0; i < len(c.Households); i++ {
		for j := 0; j < len(c.Households[i].Members); j++ {
			if c.Households[i].Members[j].IsEmployable() {
				labourforce += 1
				if !c.Households[i].Members[j].IsEmployed() {
					unemployed += 1
				}
			}
		}
	}
	c.LabourForce = labourforce
	c.Unemployed = unemployed
}
