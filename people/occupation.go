package people

import (
	"math/rand"
)

// CareerLevel defines the levels in a person's career
type CareerLevel string

const (
	Unemployed     CareerLevel = "Unemployed"
	EntryLevel     CareerLevel = "Entry Level"
	MidLevel       CareerLevel = "Mid Level"
	SeniorLevel    CareerLevel = "Senior Level"
	ExecutiveLevel CareerLevel = "Executive Level"
)

// Occupation represents a job and its salary range based on career level
type Occupation struct {
	Name        string
	SalaryRange map[CareerLevel][2]int // Min and max salary for each career level
}

// Predefined occupations and salary ranges
var occupations = []Occupation{
	{
		Name: "Software Engineer",
		SalaryRange: map[CareerLevel][2]int{
			EntryLevel:     {40000, 60000},
			MidLevel:       {60000, 90000},
			SeniorLevel:    {90000, 120000},
			ExecutiveLevel: {120000, 150000},
		},
	},
	{
		Name: "Teacher",
		SalaryRange: map[CareerLevel][2]int{
			EntryLevel:     {30000, 40000},
			MidLevel:       {40000, 60000},
			SeniorLevel:    {60000, 80000},
			ExecutiveLevel: {80000, 100000},
		},
	},
	{
		Name: "Doctor",
		SalaryRange: map[CareerLevel][2]int{
			EntryLevel:     {70000, 90000},
			MidLevel:       {90000, 130000},
			SeniorLevel:    {130000, 180000},
			ExecutiveLevel: {180000, 250000},
		},
	},
	{
		Name: "Artist",
		SalaryRange: map[CareerLevel][2]int{
			EntryLevel:     {20000, 30000},
			MidLevel:       {30000, 50000},
			SeniorLevel:    {50000, 70000},
			ExecutiveLevel: {70000, 90000},
		},
	},
}

// getRandomOccupationAndSalary randomizes occupation and salary based on age and career level
func getRandomOccupationAndSalary(age int) (string, CareerLevel, float64) {
	if age < 16 {
		return "", Unemployed, 0
	}

	// Determine career level based on age
	var careerLevel CareerLevel
	switch {
	case age < 25:
		careerLevel = EntryLevel
	case age < 35:
		careerLevel = MidLevel
	case age < 50:
		careerLevel = SeniorLevel
	case age < 70:
		careerLevel = ExecutiveLevel
	default:
		careerLevel = Unemployed
	}

	// Randomly select an occupation
	selectedOccupation := occupations[rand.Intn(len(occupations))]

	// Get salary range for the career level
	salaryRange := selectedOccupation.SalaryRange[careerLevel]
	salary := float64(rand.Intn(salaryRange[1]-salaryRange[0]+1) + salaryRange[0])

	return selectedOccupation.Name, careerLevel, salary
}
