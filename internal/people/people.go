package people

import (
	"math"
	"math/rand"
	"slices"

	"github.com/janithl/citylyf/internal/economy"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/utils"
)

const maleMeanAge = 37.0
const femaleMeanAge = 39.0
const ageStdDev = 15.0

type PeopleService struct {
	LastHouseholdID int
	LastPersonID    int
	LastFiveNames   []string
	LastTenFamilies []string
}

func NewPeopleService() *PeopleService {
	return &PeopleService{
		LastHouseholdID: 1000,  // start Household IDs from 1000
		LastPersonID:    10000, // start IDs from 10000
		LastFiveNames:   make([]string, 5),
		LastTenFamilies: make([]string, 10),
	}
}

func (ps *PeopleService) CreateRandomPerson(minAge int, maxAge int) *entities.Person {
	var name, familyName string
	var gender entities.Gender

	for name == "" || slices.Contains(ps.LastFiveNames, name) || slices.Contains(ps.LastTenFamilies, familyName) { // prevent repeating names
		name, familyName, gender = getNameAndGender()
	}
	ps.LastFiveNames = utils.AddFifo(ps.LastFiveNames, name, 5)
	ps.LastTenFamilies = utils.AddFifo(ps.LastTenFamilies, familyName, 10)

	meanAge := maleMeanAge
	if gender == entities.Female {
		meanAge = femaleMeanAge
	}
	age := getAge(meanAge, ageStdDev, minAge, maxAge)
	education := getEducationLevel(age)
	careerLevel := getCareerLevel(age, education)

	var job economy.IndustryJob
	var salary float64
	if careerLevel != entities.Unemployed {
		job, salary = economy.GetIndustryJob(education, careerLevel)
	}

	savings := salary * (float64(rand.Intn(50)) / 100) * math.Max(float64(age-25), 1)

	ps.LastPersonID += 1
	return &entities.Person{
		ID:             ps.LastPersonID,
		FirstName:      name,
		FamilyName:     familyName,
		Birthdate:      getRandomBirthdate(age),
		Gender:         gender,
		EducationLevel: education,
		Occupation:     job.Job,
		Industry:       job.Industry,
		CareerLevel:    careerLevel,
		AnnualIncome:   int(salary),
		Savings:        int(savings),
		Relationship:   entities.GetRelationshipStatus(age),
	}
}

func (ps *PeopleService) CreateHousehold() *entities.Household {
	var p, q *entities.Person
	ps.LastHouseholdID += 1
	household := &entities.Household{
		ID:         ps.LastHouseholdID,
		MemberIDs:  []int{},
		MoveInDate: entities.Sim.Date,
		LastPayDay: entities.Sim.Date,
	}

	p = ps.CreateRandomPerson(16, 100)
	entities.Sim.People.AddPerson(p)
	household.MemberIDs = append(household.MemberIDs, p.ID)
	household.Savings = p.Savings

	if p.Relationship == entities.Married {
		q = ps.CreateRandomPerson(int(math.Max(entities.AgeOfAdulthood, float64(p.Age()-15))), p.Age()+15)
		entities.Sim.People.AddPerson(q)
		q.Relationship = entities.Married
		household.Savings += q.Savings
		if rand.Intn(100) < 80 {
			q.FamilyName = p.FamilyName
		}

		household.MemberIDs = append(household.MemberIDs, q.ID)
	}

	if rand.Intn(100) < 58 {
		kidIDs := ps.createKids(p, q)
		household.MemberIDs = append(household.MemberIDs, kidIDs...)
	}

	return household
}

func (ps *PeopleService) createKids(p *entities.Person, q *entities.Person) []int {
	var kidIDs []int
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
			kidIDs = append(kidIDs, kid.ID)
		}
	}

	return kidIDs
}
