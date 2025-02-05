package economy

import (
	"math/rand"
	"slices"

	"github.com/janithl/citylyf/internal/entities"
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
		Industry:        entities.Technology,
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
		Industry:        entities.Technology,
		Job:             entities.QualityEngineer,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {40000, 60000},
			entities.MidLevel:       {60000, 90000},
			entities.SeniorLevel:    {90000, 120000},
			entities.ExecutiveLevel: {120000, 150000},
		},
	},
	{
		Industry:        entities.Telecommunications,
		Job:             entities.NetworkEngineer,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {60000, 75000},
			entities.MidLevel:       {75000, 100000},
			entities.SeniorLevel:    {100000, 125000},
			entities.ExecutiveLevel: {125000, 150000},
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
		Industry:        entities.Healthcare,
		Job:             entities.Nurse,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {40000, 60000},
			entities.MidLevel:       {60000, 100000},
			entities.SeniorLevel:    {100000, 150000},
			entities.ExecutiveLevel: {150000, 200000},
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
	{
		Industry:        entities.Technology,
		Job:             entities.CybersecurityAnalyst,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {55000, 80000},
			entities.MidLevel:       {80000, 120000},
			entities.SeniorLevel:    {120000, 160000},
			entities.ExecutiveLevel: {160000, 200000},
		},
	},
	{
		Industry:        entities.Healthcare,
		Job:             entities.Paramedic,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {40000, 60000},
			entities.MidLevel:       {60000, 85000},
			entities.SeniorLevel:    {85000, 110000},
			entities.ExecutiveLevel: {110000, 130000},
		},
	},
	{
		Industry:        entities.Retail,
		Job:             entities.SupplyChainManager,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {50000, 75000},
			entities.MidLevel:       {75000, 110000},
			entities.SeniorLevel:    {110000, 140000},
			entities.ExecutiveLevel: {140000, 180000},
		},
	},
	{
		Industry:        entities.Energy,
		Job:             entities.Geologist,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {60000, 85000},
			entities.MidLevel:       {85000, 120000},
			entities.SeniorLevel:    {120000, 160000},
			entities.ExecutiveLevel: {160000, 210000},
		},
	},
	{
		Industry:        entities.Finance,
		Job:             entities.AIResearcher,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {70000, 100000},
			entities.MidLevel:       {100000, 140000},
			entities.SeniorLevel:    {140000, 190000},
			entities.ExecutiveLevel: {190000, 250000},
		},
	},
	{
		Industry:        entities.Automobile,
		Job:             entities.MechanicalEngineer,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {45000, 65000},
			entities.MidLevel:       {65000, 90000},
			entities.SeniorLevel:    {90000, 120000},
			entities.ExecutiveLevel: {120000, 160000},
		},
	},
	{
		Industry:        entities.Retail,
		Job:             entities.StoreManager,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {30000, 45000},
			entities.MidLevel:       {45000, 65000},
			entities.SeniorLevel:    {65000, 90000},
			entities.ExecutiveLevel: {90000, 120000},
		},
	},
	{
		Industry:        entities.Finance,
		Job:             entities.FinancialAnalyst,
		EducationLevels: []entities.EducationLevel{entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {50000, 80000},
			entities.MidLevel:       {80000, 120000},
			entities.SeniorLevel:    {120000, 160000},
			entities.ExecutiveLevel: {160000, 220000},
		},
	},
	{
		Industry:        entities.Energy,
		Job:             entities.ElectricalEngineer,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {55000, 75000},
			entities.MidLevel:       {75000, 100000},
			entities.SeniorLevel:    {100000, 140000},
			entities.ExecutiveLevel: {140000, 190000},
		},
	},
	{
		Industry:        entities.Agriculture,
		Job:             entities.FarmManager,
		EducationLevels: []entities.EducationLevel{entities.Unqualified, entities.HighSchool, entities.University},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {30000, 45000},
			entities.MidLevel:       {45000, 65000},
			entities.SeniorLevel:    {65000, 90000},
			entities.ExecutiveLevel: {90000, 120000},
		},
	},
	{
		Industry:        entities.Construction,
		Job:             entities.CivilEngineer,
		EducationLevels: []entities.EducationLevel{entities.HighSchool, entities.University, entities.Postgrad},
		SalaryRange: map[entities.CareerLevel][2]int{
			entities.EntryLevel:     {40000, 60000},
			entities.MidLevel:       {60000, 90000},
			entities.SeniorLevel:    {90000, 120000},
			entities.ExecutiveLevel: {120000, 150000},
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
