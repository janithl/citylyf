package people

import (
	"math/rand"
	"slices"
)

// EducationLevel defines the levels in a person's education
type EducationLevel string

const (
	Unqualified EducationLevel = "Unqualified"
	HighSchool  EducationLevel = "High School"
	University  EducationLevel = "University"
	Postgrad    EducationLevel = "Postgrad"
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
	Name            string
	EducationLevels []EducationLevel
	SalaryRange     map[CareerLevel][2]int // Min and max salary for each career level
}

// Predefined occupations and salary ranges
var occupations = []Occupation{
	{
		Name:            "Software Engineer",
		EducationLevels: []EducationLevel{Unqualified, HighSchool, University, Postgrad},
		SalaryRange: map[CareerLevel][2]int{
			EntryLevel:     {40000, 60000},
			MidLevel:       {60000, 90000},
			SeniorLevel:    {90000, 120000},
			ExecutiveLevel: {120000, 150000},
		},
	},
	{
		Name:            "Teacher",
		EducationLevels: []EducationLevel{University, Postgrad},
		SalaryRange: map[CareerLevel][2]int{
			EntryLevel:     {30000, 40000},
			MidLevel:       {40000, 60000},
			SeniorLevel:    {60000, 80000},
			ExecutiveLevel: {80000, 100000},
		},
	},
	{
		Name:            "Doctor",
		EducationLevels: []EducationLevel{University, Postgrad},
		SalaryRange: map[CareerLevel][2]int{
			EntryLevel:     {70000, 90000},
			MidLevel:       {90000, 130000},
			SeniorLevel:    {130000, 180000},
			ExecutiveLevel: {180000, 250000},
		},
	},
	{
		Name:            "Artist",
		EducationLevels: []EducationLevel{Unqualified, HighSchool, University, Postgrad},
		SalaryRange: map[CareerLevel][2]int{
			EntryLevel:     {20000, 30000},
			MidLevel:       {30000, 50000},
			SeniorLevel:    {50000, 70000},
			ExecutiveLevel: {70000, 90000},
		},
	},
}

// getEducationLevel returns education level based on age
func getEducationLevel(age int) EducationLevel {
	if age <= AgeOfAdulthood {
		return Unqualified
	}

	randomEducation := rand.Intn(100)
	if age < 24 {
		if randomEducation < 60 {
			return HighSchool
		} else {
			return Unqualified
		}

	}

	if randomEducation < 2 {
		return Postgrad
	} else if randomEducation < 30 {
		return University
	} else if randomEducation < 70 {
		return HighSchool
	} else {
		return Unqualified
	}
}

// getRandomOccupationAndSalary randomizes occupation and salary based on age and career level
func getRandomOccupationAndSalary(age int, education EducationLevel) (string, CareerLevel, float64) {
	if age < AgeOfAdulthood || age > ageOfRetirement {
		return "", Unemployed, 0
	}

	// Determine career level based on age
	var careerLevel CareerLevel
	switch education {
	case Unqualified, HighSchool:
		switch {
		case age < 20:
			careerLevel = EntryLevel
		case age < 35:
			careerLevel = MidLevel
		case age < 50:
			careerLevel = SeniorLevel
		default:
			careerLevel = ExecutiveLevel
		}

	case University:
		switch {
		case age < 22:
			careerLevel = EntryLevel
		case age < 30:
			careerLevel = MidLevel
		case age < 45:
			careerLevel = SeniorLevel
		default:
			careerLevel = ExecutiveLevel
		}

	case Postgrad:

		switch {
		case age < 24:
			careerLevel = EntryLevel
		case age < 30:
			careerLevel = MidLevel
		case age < 40:
			careerLevel = SeniorLevel
		default:
			careerLevel = ExecutiveLevel
		}
	}

	// Randomly select an occupation
	var selectedOccupation Occupation
	for !slices.Contains(selectedOccupation.EducationLevels, education) {
		selectedOccupation = occupations[rand.Intn(len(occupations))]
	}

	// Get salary range for the career level
	salaryRange := selectedOccupation.SalaryRange[careerLevel]
	salary := float64(rand.Intn(salaryRange[1]-salaryRange[0]+1) + salaryRange[0])

	return selectedOccupation.Name, careerLevel, salary
}
