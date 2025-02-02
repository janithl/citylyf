package entities

// EducationLevel defines the levels in a person's education
type EducationLevel string

const (
	Unqualified EducationLevel = "Unqualified"
	HighSchool  EducationLevel = "High School"
	University  EducationLevel = "University"
	Postgrad    EducationLevel = "Postgrad"
)
