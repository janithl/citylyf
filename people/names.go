package people

import (
	"math/rand"
)

var familyNames = []string{
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

var namesMale = []string{
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

var namesFemale = []string{
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

func getNameAndGender() (string, string, Gender) {
	randomGender := rand.Intn(100)
	gender := Other
	switch {
	case randomGender < 49:
		gender = Male
	case randomGender < 98:
		gender = Female
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
