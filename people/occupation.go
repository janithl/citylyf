package people

import (
	"citylyf/entities"
	"math/rand"
)

// getEducationLevel returns education level based on age
func getEducationLevel(age int) entities.EducationLevel {
	if age <= entities.AgeOfAdulthood {
		return entities.Unqualified
	}

	randomEducation := rand.Intn(100)
	if age < 24 {
		if randomEducation < 60 {
			return entities.HighSchool
		} else {
			return entities.Unqualified
		}

	}

	if randomEducation < 2 {
		return entities.Postgrad
	} else if randomEducation < 30 {
		return entities.University
	} else if randomEducation < 70 {
		return entities.HighSchool
	} else {
		return entities.Unqualified
	}
}

// getCareerLevel returns career level based on age and education
func getCareerLevel(age int, education entities.EducationLevel) entities.CareerLevel {
	if age < entities.AgeOfAdulthood || age > entities.AgeOfRetirement {
		return entities.Unemployed
	}

	// Determine career level based on age
	var careerLevel entities.CareerLevel
	switch education {
	case entities.Unqualified, entities.HighSchool:
		switch {
		case age < 20:
			careerLevel = entities.EntryLevel
		case age < 35:
			careerLevel = entities.MidLevel
		case age < 50:
			careerLevel = entities.SeniorLevel
		default:
			careerLevel = entities.ExecutiveLevel
		}

	case entities.University:
		switch {
		case age < 22:
			careerLevel = entities.EntryLevel
		case age < 30:
			careerLevel = entities.MidLevel
		case age < 45:
			careerLevel = entities.SeniorLevel
		default:
			careerLevel = entities.ExecutiveLevel
		}

	case entities.Postgrad:
		switch {
		case age < 24:
			careerLevel = entities.EntryLevel
		case age < 30:
			careerLevel = entities.MidLevel
		case age < 40:
			careerLevel = entities.SeniorLevel
		default:
			careerLevel = entities.ExecutiveLevel
		}
	}

	return careerLevel
}
