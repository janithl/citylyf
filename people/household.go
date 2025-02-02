package people

import (
	"math"
	"math/rand"
)

type Household struct {
	Members []Person // Family members
}

func (h *Household) FamilyName() string {
	return h.Members[0].FamilyName
}

func (h *Household) AnnualIncome() int {
	income := 0
	for i := 0; i < len(h.Members); i++ {
		income += h.Members[i].AnnualIncome
	}
	return income
}
func (h *Household) Wealth() int {
	wealth := 0
	for i := 0; i < len(h.Members); i++ {
		wealth += h.Members[i].Wealth
	}
	return wealth
}

func CreateHousehold() Household {
	var p, q Person
	var household []Person

	p = createRandomPerson(16, 100)
	household = append(household, p)

	if p.Relationship == Married {
		q = createRandomPerson(int(math.Max(AgeOfAdulthood, float64(p.Age()-15))), p.Age()+15)
		q.Relationship = Married
		if rand.Intn(100) < 80 {
			q.FamilyName = p.FamilyName
		}

		household = append(household, q)
	}

	if rand.Intn(100) < 58 {
		kids := createKids(p, q)
		household = append(household, kids...)
	}

	return Household{
		Members: household,
	}
}

func createKids(p Person, q Person) []Person {
	var kids []Person
	numberOfKids := 0

	randomKids := rand.Intn(100)
	switch {
	case randomKids < 34:
		numberOfKids = 0
	case randomKids < 47:
		numberOfKids = 1
	case randomKids < 72:
		numberOfKids = 2
	case randomKids < 86:
		numberOfKids = 3
	case randomKids < 92:
		numberOfKids = 4
	default:
		numberOfKids = 5
	}

	for i := 0; i < numberOfKids; i++ {
		var kid Person

		if p.Relationship == Married {
			if q.Age() == 0 {
				kid = createRandomPerson(0, p.Age()-AgeOfAdulthood)
			} else {
				parentMaxAge := p.Age()
				if q.Age() > p.Age() {
					parentMaxAge = q.Age()
				}
				kid = createRandomPerson(0, parentMaxAge-AgeOfAdulthood)
			}
		}

		if p.Relationship == Divorced || p.Relationship == Widowed {
			kid = createRandomPerson(5, p.Age()-AgeOfAdulthood)
		}

		if kid.FirstName != "" {
			kid.Relationship = Single
			kid.FamilyName = p.FamilyName
			kids = append(kids, kid)
		}
	}

	return kids
}
