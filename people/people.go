package people

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Person struct {
	ID           int         // Unique identifier for the person
	FirstName    string      // Name of the person
	FamilyName   string      // Family name
	Birthdate    time.Time   // Age of the person
	Gender       Gender      // Gender of the person
	Occupation   string      // Job title or role
	CareerLevel  CareerLevel // Their career level
	AnnualIncome int         // Annual income
	Wealth       int         // Total assets or savings
	Relationship RelationshipStatus
}

func (p *Person) Age() int {
	currentDate := time.Now()
	duration := currentDate.Sub(p.Birthdate)
	hoursPerYear := 24 * 365.25
	return int(duration.Hours() / hoursPerYear)
}

func (p *Person) String() string {
	return fmt.Sprintf("%20s %20s  %3d %6s %10s %20s %20s %10d/yearly", p.FirstName, p.FamilyName, p.Age(), p.Gender, p.Relationship, p.CareerLevel, p.Occupation, p.AnnualIncome)
}

func createRandomPerson(minAge int, maxAge int) Person {
	name, familyName, gender := getNameAndGender()
	meanAge := 37.0
	if gender == Female {
		meanAge = 39.0
	}
	age := getAge(meanAge, 15, minAge, maxAge)
	occupation, careerLevel, salary := getRandomOccupationAndSalary(age)
	savings := salary * (float64(rand.Intn(50)) / 100) * math.Max(float64(age-25), 1)

	return Person{
		ID:           rand.Intn(9999) + 10000,
		FirstName:    name,
		FamilyName:   familyName,
		Birthdate:    getRandomBirthdate(age),
		Gender:       gender,
		Occupation:   occupation,
		CareerLevel:  careerLevel,
		AnnualIncome: int(salary),
		Wealth:       int(savings),
		Relationship: getRelationshipStatus(age),
	}
}
