package main

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

// getRandomOccupationAndSalary randomizes occupation and salary based on age and career level
func getRandomOccupationAndSalary(age int) (string, CareerLevel, float64) {
	if age < 16 {
		return "", Unemployed, 0
	}

	// Predefined occupations and salary ranges
	occupations := []Occupation{
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

// Gender defines the person's gender
type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
	Other  Gender = "Other"
)

// RelationshipStatus defines if the person is married etc.
type RelationshipStatus string

const (
	Single   RelationshipStatus = "Single"
	Married  RelationshipStatus = "Married"
	Divorced RelationshipStatus = "Divorced"
	Widowed  RelationshipStatus = "Widowed"
)

func getNameAndGender() (string, string, Gender) {
	randomGender := rand.Intn(100)
	gender := Other
	switch {
	case randomGender < 49:
		gender = Male
	case randomGender < 98:
		gender = Female
	}

	familyNames := []string{
		"Maleck",
		"Tallis",
		"Ormston",
		"Bruhnicke",
		"Mathely",
		"Darbishire",
		"O'Hagirtie",
		"Zarb",
		"Gillison",
		"Karpenya",
		"Bonds",
		"Cowherd",
		"Penreth",
		"Drysdell",
		"Bohea",
		"Shermar",
		"Pennicott",
		"Pickless",
		"Frow",
		"Ugoletti",
		"Tunnicliffe",
		"Ruthen",
		"Colquite",
		"Hammerberger",
		"Van Bruggen",
		"Ledrun",
		"McMurdo",
		"Chellam",
		"Claypool",
		"Whittier",
		"Callard",
		"Southorn",
		"Sprankling",
		"Hutchason",
		"Enrich",
		"Matej",
		"Campsall",
		"Gerardi",
		"Solomonides",
		"Wimes",
		"Josephsen",
		"Abazi",
		"MacKibbon",
		"Stanway",
		"Skeleton",
		"Pavyer",
		"Breznovic",
		"Jerzyk",
		"Goding",
		"Groneway",
	}

	namesMale := []string{
		"Gregorius",
		"Talbert",
		"Esteban",
		"Dudley",
		"Randolf",
		"Judon",
		"Englebert",
		"Maximo",
		"Jeremie",
		"Gabriele",
		"Karlan",
		"Gill",
		"Tyson",
		"Tymothy",
		"Cheston",
		"Putnam",
		"Oswell",
		"Justen",
		"Gerik",
		"Abe",
		"Lenci",
		"Marion",
		"Holly",
		"Benedick",
		"Winfield",
		"Fabio",
		"Arne",
		"Buddie",
		"Florian",
		"Zollie",
		"Tobias",
		"Godfry",
		"Nev",
		"Dud",
		"Jae",
		"Purcell",
		"Fabian",
		"Aldon",
		"Elisha",
		"Whitney",
		"Elwood",
		"Killy",
		"Skipp",
		"Grenville",
		"Jeffry",
		"Raymund",
		"Bil",
		"Shelden",
		"Dun",
		"Laird",
	}

	namesFemale := []string{
		"Shawn",
		"Kaia",
		"Betta",
		"Phylys",
		"Correna",
		"Rosanne",
		"Vivienne",
		"Debby",
		"Anselma",
		"Sydney",
		"Annabal",
		"Rozella",
		"Cassey",
		"Vania",
		"Bamby",
		"Gracia",
		"Kippie",
		"Raquela",
		"Chery",
		"Myrah",
		"Brandais",
		"Mahala",
		"Holly-anne",
		"Jackelyn",
		"Gnni",
		"Elga",
		"Netti",
		"Chrissy",
		"Erda",
		"Jorie",
		"Estele",
		"Dita",
		"Glennis",
		"Marsha",
		"Rona",
		"Bree",
		"Nora",
		"Mireielle",
		"Brenn",
		"Allison",
		"Cecile",
		"Christian",
		"Shantee",
		"Paulita",
		"Persis",
		"Lilah",
		"Celina",
		"Penny",
		"Marita",
		"Koren",
	}

	namesCombined := append(namesMale, namesFemale...)

	var name string
	if gender == Male {
		name = namesMale[rand.Intn(len(namesMale))]
	} else if gender == Female {
		name = namesFemale[rand.Intn(len(namesFemale))]
	} else {
		name = namesCombined[rand.Intn(len(namesCombined))]
	}

	return name, familyNames[rand.Intn(len(familyNames))], gender
}

// getAge generates a random age based on a bell curve
func getAge(mean, stdDev float64, minAge, maxAge int) int {
	if minAge < 0 {
		minAge = 0
	}
	if maxAge < minAge {
		maxAge = minAge
	}

	maxCalculations := 100
	for i := 0; i < maxCalculations; i++ {
		// Box-Muller transform to generate normal distribution
		u1 := rand.Float64()
		u2 := rand.Float64()
		z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

		// Scale and shift to get the desired mean and standard deviation
		age := mean + z*stdDev

		// Ensure age is within bounds
		if int(age) >= minAge && int(age) <= maxAge {
			return int(math.Round(age))
		}
	}

	return int(minAge + rand.Intn(maxAge))
}

// GenerateRandomBirthdate generates a random birthdate given the age
func GenerateRandomBirthdate(age int) time.Time {
	year := time.Now().Year() - age

	// Generate a random month (1-12)
	month := time.Month(rand.Intn(12) + 1)

	// Generate a random day based on the month and year
	// Use time.Date to determine the last day of the month
	day := rand.Intn(time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()) + 1

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func getRelationshipStatus(age int) RelationshipStatus {
	if age < 18 {
		return Single
	}

	randomRelationship := rand.Intn(100)
	var status RelationshipStatus
	switch {
	case randomRelationship < 45:
		status = Married
	case randomRelationship < 85:
		status = Single
	case randomRelationship < 90 && age > 50:
		status = Widowed
	default:
		status = Divorced
	}

	return status
}

func createRandomPerson(minAge int, maxAge int) Person {
	name, familyName, gender := getNameAndGender()
	age := getAge(30.6, 15, minAge, maxAge)
	occupation, careerLevel, salary := getRandomOccupationAndSalary(age)
	savings := salary * (float64(rand.Intn(50)) / 100) * math.Max(float64(age-25), 1)

	return Person{
		ID:           rand.Intn(9999) + 10000,
		FirstName:    name,
		FamilyName:   familyName,
		Birthdate:    GenerateRandomBirthdate(age),
		Gender:       gender,
		Occupation:   occupation,
		CareerLevel:  careerLevel,
		AnnualIncome: int(salary),
		Wealth:       int(savings),
		Relationship: getRelationshipStatus(age),
	}
}

func createKids(p Person, q Person) []Person {
	var kids []Person
	numberOfKids := 0

	randomKids := rand.Intn(100)
	switch {
	case randomKids < 34:
		numberOfKids = 0
	case randomKids < 47:
		numberOfKids = 1
	case randomKids < 72:
		numberOfKids = 2
	case randomKids < 86:
		numberOfKids = 3
	case randomKids < 92:
		numberOfKids = 4
	default:
		numberOfKids = 5
	}

	for i := 0; i < numberOfKids; i++ {
		var kid Person

		if p.Relationship == Married && p.Age() >= 20 {
			if q.Age() >= 20 {
				kid = createRandomPerson(0, int(math.Min(float64(p.Age()), float64(q.Age())))-19)
			}
			if q.Age() == 0 {
				kid = createRandomPerson(0, p.Age()-19)
			}
		}

		if p.Relationship == Divorced || p.Relationship == Widowed {
			kid = createRandomPerson(5, p.Age()-19)
		}

		if kid.FirstName != "" {
			kid.Relationship = Single
			kid.FamilyName = p.FamilyName
			kids = append(kids, kid)
		}
	}

	return kids
}

type Household struct {
	Members []Person // Family members
}

func (h *Household) FamilyName() string {
	return h.Members[0].FamilyName
}

func createHousehold() Household {
	var p, q Person
	var household []Person

	p = createRandomPerson(16, 100)
	household = append(household, p)

	if p.Relationship == Married {
		q = createRandomPerson(int(math.Max(18, float64(p.Age()-15))), p.Age()+15)
		q.Relationship = Married
		if rand.Intn(100) < 80 {
			q.FamilyName = p.FamilyName
		}

		household = append(household, q)
	}

	kids := createKids(p, q)
	household = append(household, kids...)

	return Household{
		Members: household,
	}
}

func main() {
	for i := 0; i < 100; i++ {
		h := createHousehold()
		fmt.Printf("[The %s Family]:\n", h.FamilyName())
		for j := 0; j < len(h.Members); j++ {
			fmt.Printf(" |-> %s\n", h.Members[j].String())
		}
	}
}
