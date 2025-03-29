package people

import (
	"math"
	"math/rand/v2"
	"strings"

	"github.com/janithl/citylyf/internal/economy"
	"github.com/janithl/citylyf/internal/entities"
)

func CreateRandomPerson(minAge int, maxAge int) *entities.Person {
	gender := entities.GetRandomGender()
	name, familyName := entities.Sim.NameService.GetPersonName(gender)

	meanAge := entities.MeanAgeMale
	if gender == entities.Female {
		meanAge = entities.MeanAgeFemale
	}

	ageY, ageM := getAge(meanAge, entities.AgeStdDev, minAge, maxAge)
	education := getEducationLevel(ageY)
	careerLevel := getCareerLevel(ageY, education)

	var job economy.IndustryJob
	var salary float64
	if careerLevel != entities.Unemployed {
		job, salary = economy.GetIndustryJob(education, careerLevel)
	}

	savings := salary * rand.Float64() * 0.5 * math.Max(float64(ageY-25), 1)

	return &entities.Person{
		FirstName:      name,
		FamilyName:     familyName,
		Birthdate:      getRandomBirthdate(ageY, ageM),
		Gender:         gender,
		EducationLevel: education,
		Occupation:     job.Job,
		Industry:       job.Industry,
		CareerLevel:    careerLevel,
		AnnualIncome:   int(salary),
		Savings:        int(savings),
		Relationship:   entities.GetRelationshipStatus(ageY),
	}
}

func CreateHousehold() *entities.Household {
	var p, q *entities.Person
	householdID := entities.Sim.GetNextID()

	household := &entities.Household{
		ID:         householdID,
		MemberIDs:  []int{},
		MoveInDate: entities.Sim.Date,
		LastPayDay: entities.Sim.Date,
	}

	p = CreateRandomPerson(16, 100)
	p.ID = entities.Sim.GetNextID()
	entities.Sim.People.AddPerson(p)
	household.MemberIDs = append(household.MemberIDs, p.ID)
	household.Savings = p.Savings

	if p.Relationship == entities.Married {
		q = CreateRandomPerson(int(math.Max(entities.AgeOfAdulthood, float64(p.Age()-15))), p.Age()+15)
		q.ID = entities.Sim.GetNextID()
		entities.Sim.People.AddPerson(q)
		q.Relationship = entities.Married
		household.Savings += q.Savings
		if rand.IntN(100) < 80 {
			q.FamilyName = p.FamilyName
		}

		household.MemberIDs = append(household.MemberIDs, q.ID)
	}

	if rand.IntN(100) < 58 {
		kids := createKids(p, q, getNumberOfKids())
		for _, kid := range kids {
			kid.ID = entities.Sim.GetNextID()
			entities.Sim.People.AddPerson(kid)
			household.MemberIDs = append(household.MemberIDs, kid.ID)
		}
	}

	return household
}

// RemoveHousehold removes a household and its members from the Sim, and removes them from their jobs
func RemoveHousehold(household *entities.Household) {
	for _, memberID := range household.MemberIDs {
		member := entities.Sim.People.GetPerson(memberID)
		if member != nil {
			entities.Sim.Companies.RemoveEmployeeFromTheirCompany(member)
			entities.Sim.People.RemovePerson(memberID)
		}
	}
	delete(entities.Sim.People.Households, household.ID)
}

func getNumberOfKids() int {
	randomKids := rand.IntN(100)
	switch {
	case randomKids < 34:
		return 0
	case randomKids < 47:
		return 1
	case randomKids < 72:
		return 2
	case randomKids < 86:
		return 3
	case randomKids < 92:
		return 4
	default:
		return 5
	}
}

func createKids(p *entities.Person, q *entities.Person, numberOfKids int) []*entities.Person {
	var kids []*entities.Person

	if numberOfKids == 0 {
		return kids
	}

	familyName := p.FamilyName
	if q != nil && !strings.Contains(familyName, "-") && rand.Float32() < 0.1 { // 10% of surnames are double‑barrelled
		familyName += "-" + q.FamilyName
	}

	for len(kids) < numberOfKids {
		parentMaxAge := p.Age()
		kidMinAge := 0
		if p.Relationship == entities.Married && q != nil && q.Age() > parentMaxAge {
			parentMaxAge = q.Age()
		} else if p.Relationship == entities.Widowed || p.Relationship == entities.Divorced {
			kidMinAge = 1
		}

		if kid := CreateRandomPerson(kidMinAge, parentMaxAge-entities.AgeOfAdulthood); kid != nil {
			kid.Relationship = entities.Single
			kid.FamilyName = familyName
			kids = append(kids, kid)
		}
	}

	return kids
}
