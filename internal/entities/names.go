package entities

import (
	"math/rand/v2"
	"slices"

	"github.com/janithl/citylyf/internal/ui/assets"
	"github.com/janithl/citylyf/internal/utils"
)

type NameService struct {
	LastTenNames     []string
	LastTenFamilies  []string
	LastTenCompanies []string
	LastTenPlaces    []string
}

func (ns *NameService) getFamilyName() string {
	familyName := ""
	for familyName == "" || slices.Contains(ns.LastTenFamilies, familyName) { // prevent repeating names
		familyNames, exists := assets.Assets.Names["familyNames"]
		if exists {
			familyName = familyNames[rand.IntN(len(familyNames))]
		}
	}
	ns.LastTenFamilies = utils.AddFifo(ns.LastTenFamilies, familyName, 10)
	return familyName
}

func (ns *NameService) GetPlaceName() string {
	placeName := ""
	placeNames, exists := assets.Assets.Names["placeNames"]
	for placeName == "" || slices.Contains(ns.LastTenPlaces, placeName) { // prevent repeating names
		if exists {
			placeName = placeNames[rand.IntN(len(placeNames))]
		} else {
			return ""
		}
	}
	ns.LastTenPlaces = utils.AddFifo(ns.LastTenPlaces, placeName, 10)
	return placeName
}

func (ns *NameService) GetRoadName() string {
	roadName := ""
	suffix := ""

	roadSuffixes, exists := assets.Assets.Names["roadSuffixes"]
	if exists {
		suffix = roadSuffixes[rand.IntN(len(roadSuffixes))]
	}

	randomNumber := rand.Float32()
	switch {
	case randomNumber < 0.7: // 70% of roads have place names
		roadName = ns.GetPlaceName()
	default: // the rest are family names
		roadName = ns.getFamilyName()
	}

	return roadName + " " + suffix
}

func (ns *NameService) GetPersonName(gender Gender) (string, string) {
	firstName := ""
	for firstName == "" || slices.Contains(ns.LastTenNames, firstName) { // prevent repeating names
		switch gender {
		case Male:
			maleNames, exists := assets.Assets.Names["maleNames"]
			if exists {
				firstName = maleNames[rand.IntN(len(maleNames))]
			}
		case Female:
			femaleNames, exists := assets.Assets.Names["femaleNames"]
			if exists {
				firstName = femaleNames[rand.IntN(len(femaleNames))]
			}
		default:
			otherNames, exists := assets.Assets.Names["otherNames"]
			if exists {
				firstName = otherNames[rand.IntN(len(otherNames))]
			}
		}
	}

	ns.LastTenNames = utils.AddFifo(ns.LastTenNames, firstName, 10)

	if rand.Float32() < 0.1 { // 10% of surnames are double‑barrelled
		return firstName, ns.getFamilyName() + "-" + ns.getFamilyName()
	}

	return firstName, ns.getFamilyName()
}

func (ns *NameService) GetCompanyName() string {
	companyName := ""
	suffix := ""

	companySuffixes, exists := assets.Assets.Names["companySuffixes"]
	if exists {
		suffix = companySuffixes[rand.IntN(len(companySuffixes))]
	} else {
		return ""
	}

	companyNames, exists := assets.Assets.Names["companyNames"]
	for companyName == "" || slices.Contains(ns.LastTenCompanies, companyName) { // prevent repeating names
		randomNumber := rand.Float32()
		switch {
		case randomNumber < 0.65: // 65% are generic company names
			if exists {
				companyName = companyNames[rand.IntN(len(companyNames))]
			} else {
				return ""
			}
		case randomNumber < 0.85: // the next 20% are family names
			companyName = ns.getFamilyName()
		default: // the rest are place names
			companyName = ns.GetPlaceName()
		}
	}

	ns.LastTenCompanies = utils.AddFifo(ns.LastTenCompanies, companyName, 10)
	return companyName + " " + suffix
}

func NewNameService() *NameService {
	return &NameService{
		LastTenNames:     make([]string, 10),
		LastTenFamilies:  make([]string, 10),
		LastTenCompanies: make([]string, 10),
		LastTenPlaces:    make([]string, 10),
	}
}
