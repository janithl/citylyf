package entities

import (
	"fmt"
	"time"
)

const (
	AgeOfAdulthood  = 18
	AgeOfRetirement = 70
	HoursPerDay     = 24
	DaysPerYear     = 365.25
)

const HoursPerYear = HoursPerDay * DaysPerYear

type Person struct {
	ID             int            // Unique identifier for the person
	FirstName      string         // Name of the person
	FamilyName     string         // Family name
	Birthdate      time.Time      // Age of the person
	Gender         Gender         // Gender of the person
	EducationLevel EducationLevel // Person's education level
	Occupation     Job            // Job title or role
	Industry       Industry       // Industry of the job
	CareerLevel    CareerLevel    // Their career level
	EmployerID     int            // Their employer id
	AnnualIncome   int            // Annual income
	Savings        int            // Total personal savings
	Relationship   RelationshipStatus
}

func (p *Person) Age() int {
	duration := Sim.Date.Sub(p.Birthdate)
	return int(duration.Hours() / HoursPerYear)
}

func (p *Person) IsEmployable() bool {
	return p.Age() > AgeOfAdulthood || p.CareerLevel != Unemployed
}

func (p *Person) IsEmployed() bool {
	return p.EmployerID != 0
}

func (p *Person) String() string {
	return fmt.Sprintf("%20s %20s  %3d (%4d) %6s %10s %15s %20s %25s %5d %10d/yearly", p.FirstName, p.FamilyName, p.Age(), p.Birthdate.Year(), p.Gender, p.Relationship, p.EducationLevel, p.CareerLevel, p.Occupation, p.EmployerID, p.AnnualIncome)
}
