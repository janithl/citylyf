package economy

import (
	"citylyf/entities"
	"math/rand"
	"slices"
)

// IndustryJob represents the jobs in an industry
type IndustryJob struct {
	Industry        entities.Industry
	Job             entities.Job
	EducationLevels []entities.EducationLevel
	SalaryRange     map[entities.CareerLevel][2]int // Min and max salary for each career level
}

// Predefined occupations and salary ranges
var Jobs = []IndustryJob{
	{
		Industry:        entities.Software,
		Job:             entities.SoftwareEngineer,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {40000, 60000},
			entities.MidLevel:       {60000, 90000},
			entities.SeniorLevel:    {90000, 120000},
			entities.ExecutiveLevel: {120000, 150000},
		},
	},
	{
		Industry:        entities.Education,
		Job:             entities.Teacher,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {30000, 40000},
			entities.MidLevel:       {40000, 60000},
			entities.SeniorLevel:    {60000, 80000},
			entities.ExecutiveLevel: {80000, 100000},
		},
	},
	{
		Industry:        entities.Healthcare,
		Job:             entities.Doctor,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {70000, 90000},
			entities.MidLevel:       {90000, 130000},
			entities.SeniorLevel:    {130000, 180000},
			entities.ExecutiveLevel: {180000, 250000},
		},
	},
	{
		Industry:        entities.Creative,
		Job:             entities.Artist,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {20000, 30000},
			entities.MidLevel:       {30000, 50000},
			entities.SeniorLevel:    {50000, 70000},
			entities.ExecutiveLevel: {70000, 90000},
		},
	},
}

// Randomly assigns an industry job
func GetIndustryJob(education entities.EducationLevel, careerLevel entities.CareerLevel) (IndustryJob, float64) {
	var selectedJob IndustryJob

	// Get a job suitable for the education level
	for !slices.Contains(selectedJob.EducationLevels, education) {
		selectedJob = Jobs[rand.Intn(len(Jobs))]
	}

	// Get salary range for the career level
	salaryRange := selectedJob.SalaryRange[careerLevel]
	salary := float64(rand.Intn(salaryRange[1]-salaryRange[0]+1) + salaryRange[0])

	return selectedJob, salary
}
