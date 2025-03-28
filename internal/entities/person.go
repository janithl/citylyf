package entities

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	AgeOfAdulthood  = 18
	AgeOfRetirement = 70

	MeanMarriageAge          = 30.0 // Average age for marriage
	StdDevMarriageAge        = 7.0  // Standard deviation of marriage age
	MaxMarriageAgeDifference = 15
	ProbabilityOfMarriage    = 0.009 // Annual marriage rate is about 9 marriages per 1000 people

	HoursPerDay = 24
	DaysPerYear = 365.25
)

const HoursPerYear = HoursPerDay * DaysPerYear

type Person struct {
	ID, EmployerID        int // Unique identifier for the person, and their employer
	FirstName, FamilyName string
	Birthdate             time.Time      // Age of the person
	Gender                Gender         // Gender of the person
	EducationLevel        EducationLevel // Person's education level
	Occupation            Job            // Job title or role
	Industry              Industry       // Industry of the job
	CareerLevel           CareerLevel    // Their career level
	AnnualIncome, Savings int            // Annual income and total personal savings
	Relationship          RelationshipStatus
}

func (p *Person) Age() int {
	duration := Sim.Date.Sub(p.Birthdate)
	return int(duration.Hours() / HoursPerYear)
}

func (p *Person) IsEmployable() bool {
	return p.Age() >= AgeOfAdulthood && p.CareerLevel != Retired
}

func (p *Person) IsEmployed() bool {
	return p.EmployerID != 0
}

func (p *Person) CurrentIncome() int {
	if p.IsEmployed() {
		return p.AnnualIncome
	}
	return 0
}

func (p *Person) ConsiderRetirement(removeEmployeeFromCompany func(companyID int, employeeID int)) bool {
	if p.CareerLevel != Retired &&
		p.Age() >= AgeOfRetirement &&
		rand.Intn(100) < 1+p.Age()-AgeOfRetirement { // chance of retirement if over retirement age, starts at 1% and goes up

		removeEmployeeFromCompany(p.EmployerID, p.ID)
		p.EmployerID = 0
		p.CareerLevel = Retired
		return true
	}
	return false
}

func (p *Person) String() string {
	return fmt.Sprintf("%-20s%-20s%3d (%4d) %6s %10s %15s %20s %25s %5d %10d/yearly", p.FirstName, p.FamilyName, p.Age(), p.Birthdate.Year(), p.Gender, p.Relationship, p.EducationLevel, p.CareerLevel, p.Occupation, p.EmployerID, p.AnnualIncome)
}
