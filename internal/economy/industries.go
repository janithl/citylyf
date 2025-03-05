package economy

import (
	"math"
	"math/rand"
	"slices"

	"github.com/janithl/citylyf/internal/entities"
)

// Predefined occupations and salary ranges
var Jobs = []entities.CompanyJob{
	{
		Industry:        entities.Technology,
		Job:             entities.SoftwareEngineer,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University, entities.Postgrad},
		BaseSalary:      40000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Technology,
		Job:             entities.QualityEngineer,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University, entities.Postgrad},
		BaseSalary:      40000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Telecommunications,
		Job:             entities.NetworkEngineer,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		BaseSalary:      60000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Education,
		Job:             entities.Teacher,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		BaseSalary:      30000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Healthcare,
		Job:             entities.Doctor,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		BaseSalary:      70000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Healthcare,
		Job:             entities.Nurse,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		BaseSalary:      40000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Creative,
		Job:             entities.Artist,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University, entities.Postgrad},
		BaseSalary:      20000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Technology,
		Job:             entities.CybersecurityAnalyst,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University, entities.Postgrad},
		BaseSalary:      55000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Healthcare,
		Job:             entities.Paramedic,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University},
		BaseSalary:      40000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Retail,
		Job:             entities.SupplyChainManager,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University, entities.Postgrad},
		BaseSalary:      50000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Energy,
		Job:             entities.Geologist,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		BaseSalary:      60000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Finance,
		Job:             entities.AIResearcher,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		BaseSalary:      70000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Automobile,
		Job:             entities.MechanicalEngineer,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University},
		BaseSalary:      45000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Retail,
		Job:             entities.StoreManager,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University},
		BaseSalary:      30000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Finance,
		Job:             entities.FinancialAnalyst,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		BaseSalary:      50000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Energy,
		Job:             entities.ElectricalEngineer,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University, entities.Postgrad},
		BaseSalary:      55000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Agriculture,
		Job:             entities.FarmManager,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University},
		BaseSalary:      30000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
	{
		Industry:        entities.Construction,
		Job:             entities.CivilEngineer,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University, entities.Postgrad},
		BaseSalary:      40000,
		SalaryRange:     0.5,
		Increment:       0.5,
	},
}

func GetCompanyJobOpenings(industry entities.Industry, openings int) []*entities.CompanyJob {
	jobOpenings := []*entities.CompanyJob{}
	suitableJobs := []*entities.CompanyJob{}
	for _, job := range Jobs {
		if job.Industry == industry {
			suitableJobs = append(suitableJobs, &job)
		}
	}
	for range openings {
		jobOpenings = append(jobOpenings, suitableJobs[rand.Intn(len(suitableJobs))])
	}
	return jobOpenings
}

// Randomly assigns an industry job
func GetIndustryJob(education entities.EducationLevel, careerLevel entities.CareerLevel) (entities.CompanyJob, float64) {
	var selectedJob entities.CompanyJob

	// Get a job suitable for the education level
	for !slices.Contains(selectedJob.EducationLevels, education) {
		selectedJob = Jobs[rand.Intn(len(Jobs))]
	}

	level := 0
	switch careerLevel {
	case entities.MidLevel:
		level = 1
	case entities.SeniorLevel:
		level = 2
	case entities.ExecutiveLevel:
		level = 3
	}

	// Get salary range for the career level
	salary := float64(selectedJob.BaseSalary) +
		float64(selectedJob.BaseSalary)*selectedJob.SalaryRange*rand.Float64() +
		float64(selectedJob.Increment)*math.Pow(selectedJob.Increment, float64(level))

	return selectedJob, salary
}
