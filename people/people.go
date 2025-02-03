package people

import (
	"math"
	"math/rand"

	"github.com/janithl/citylyf/economy"
	"github.com/janithl/citylyf/entities"
)

func createRandomPerson(minAge int, maxAge int) entities.Person {
	name, familyName, gender := getNameAndGender()
	meanAge := 37.0
	if gender == entities.Female {
		meanAge = 39.0
	}
	age := getAge(meanAge, 15, minAge, maxAge)
	education := getEducationLevel(age)
	careerLevel := getCareerLevel(age, education)

	var job economy.IndustryJob
	var salary float64
	if careerLevel != entities.Unemployed {
		job, salary = economy.GetIndustryJob(education, careerLevel)
	}

	savings := salary * (float64(rand.Intn(50)) / 100) * math.Max(float64(age-25), 1)

	return entities.Person{
		ID:             rand.Intn(9999) + 10000,
		FirstName:      name,
		FamilyName:     familyName,
		Birthdate:      getRandomBirthdate(age),
		Gender:         gender,
		EducationLevel: education,
		Occupation:     job.Job,
		Industry:       job.Industry,
		CareerLevel:    careerLevel,
		AnnualIncome:   int(salary),
		Wealth:         int(savings),
		Relationship:   entities.GetRelationshipStatus(age),
	}
}
