package people

import (
	"math"
	"math/rand"
	"strings"

	"github.com/janithl/citylyf/internal/economy"
	"github.com/janithl/citylyf/internal/entities"
)

const maleMeanAge = 37.0
const femaleMeanAge = 39.0
const ageStdDev = 15.0

func CreateRandomPerson(minAge int, maxAge int) *entities.Person {
	gender := entities.GetRandomGender()
	name, familyName := entities.Sim.NameService.GetPersonName(gender)

	meanAge := maleMeanAge
	if gender == entities.Female {
		meanAge = femaleMeanAge
	}
	ageY, ageM := getAge(meanAge, ageStdDev, minAge, maxAge)
	education := getEducationLevel(ageY)
	careerLevel := getCareerLevel(ageY, education)

	var job economy.IndustryJob
	var salary float64
	if careerLevel != entities.Unemployed {
		job, salary = economy.GetIndustryJob(education, careerLevel)
	}

	savings := salary * (float64(rand.Intn(50)) / 100) * math.Max(float64(ageY-25), 1)

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
		if rand.Intn(100) < 80 {
			q.FamilyName = p.FamilyName
		}

		household.MemberIDs = append(household.MemberIDs, q.ID)
	}

	if rand.Intn(100) < 58 {
		kids := createKids(p, q)
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

func createKids(p *entities.Person, q *entities.Person) []*entities.Person {
	var kids []*entities.Person
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

	familyName := p.FamilyName
	if q != nil && !strings.Contains(familyName, "-") && rand.Float32() < 0.1 { // 10% of surnames are double‑barrelled
		familyName += "-" + q.FamilyName
	}
	for i := 0; i < numberOfKids; i++ {
		var kid *entities.Person

		if p.Relationship == entities.Married {
			if q.Age() == 0 {
				kid = CreateRandomPerson(0, p.Age()-entities.AgeOfAdulthood)
			} else {
				parentMaxAge := p.Age()
				if q.Age() > p.Age() {
					parentMaxAge = q.Age()
				}
				kid = CreateRandomPerson(0, parentMaxAge-entities.AgeOfAdulthood)
			}
		}

		if p.Relationship == entities.Divorced || p.Relationship == entities.Widowed {
			kid = CreateRandomPerson(5, p.Age()-entities.AgeOfAdulthood)
		}

		if kid != nil {
			kid.Relationship = entities.Single
			kid.FamilyName = familyName
			kids = append(kids, kid)
		}
	}

	return kids
}
