package people

import (
	"math"
	"math/rand"

	"github.com/janithl/citylyf/internal/entities"
)

func CreateHousehold(ps *PeopleService) entities.Household {
	var p, q *entities.Person
	household := entities.Household{
		Members:    []entities.Person{},
		MoveInDate: entities.Sim.Date,
		LastPayDay: entities.Sim.Date,
	}

	p = ps.CreateRandomPerson(16, 100)
	entities.Sim.People.AddPerson(p)
	household.Members = append(household.Members, *p)
	household.Savings = p.Savings

	if p.Relationship == entities.Married {
		q = ps.CreateRandomPerson(int(math.Max(entities.AgeOfAdulthood, float64(p.Age()-15))), p.Age()+15)
		entities.Sim.People.AddPerson(q)
		q.Relationship = entities.Married
		household.Savings += q.Savings
		if rand.Intn(100) < 80 {
			q.FamilyName = p.FamilyName
		}

		household.Members = append(household.Members, *q)
	}

	if rand.Intn(100) < 58 {
		kids := createKids(ps, p, q)
		household.Members = append(household.Members, kids...)
	}

	return household
}

func createKids(ps *PeopleService, p *entities.Person, q *entities.Person) []entities.Person {
	var kids []entities.Person
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
		var kid *entities.Person

		if p.Relationship == entities.Married {
			if q.Age() == 0 {
				kid = ps.CreateRandomPerson(0, p.Age()-entities.AgeOfAdulthood)
			} else {
				parentMaxAge := p.Age()
				if q.Age() > p.Age() {
					parentMaxAge = q.Age()
				}
				kid = ps.CreateRandomPerson(0, parentMaxAge-entities.AgeOfAdulthood)
			}
		}

		if p.Relationship == entities.Divorced || p.Relationship == entities.Widowed {
			kid = ps.CreateRandomPerson(5, p.Age()-entities.AgeOfAdulthood)
		}

		if kid != nil {
			kid.Relationship = entities.Single
			kid.FamilyName = p.FamilyName
			entities.Sim.People.AddPerson(kid)
			kids = append(kids, *kid)
		}
	}

	return kids
}
