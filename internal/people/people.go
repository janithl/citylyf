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
	LastPersonID   int
	LastThreeNames []string
}

func NewPeopleService() *PeopleService {
	return &PeopleService{
		LastPersonID:   10000, // start IDs from 10000
		LastThreeNames: make([]string, 3),
	}
}

func (ps *PeopleService) CreateRandomPerson(minAge int, maxAge int) *entities.Person {
	var name, familyName string
	var gender entities.Gender

	for name == "" || slices.Contains(ps.LastThreeNames, name) { // prevent repeating names
		name, familyName, gender = getNameAndGender()
	}
	ps.LastThreeNames = utils.AddFifo(ps.LastThreeNames, name, 3)

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
