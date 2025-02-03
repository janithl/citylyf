package entities

type CityData struct {
	Population  int
	LabourForce int
	Unemployed  int
	Households  []Household
}

func (c *CityData) UnemploymentRate() float64 {
	return 100.0 * float64(c.Unemployed) / float64(c.LabourForce)
}

func (c *CityData) MoveIn(h Household) {
	c.Households = append(c.Households, h)
	c.Population += len(h.Members)
}

func (c *CityData) CalculateUnemployment() {
	labourforce, unemployed := 0, 0
	for i := 0; i < len(c.Households); i++ {
		for j := 0; j < len(c.Households[i].Members); j++ {
			if c.Households[i].Members[j].Age() > AgeOfAdulthood || c.Households[i].Members[j].CareerLevel != Unemployed {
				labourforce += 1
				if c.Households[i].Members[j].EmployerID == 0 {
					unemployed += 1
				}
			}
		}
	}
	c.LabourForce = labourforce
	c.Unemployed = unemployed
}
