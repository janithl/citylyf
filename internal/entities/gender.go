package entities

import "math/rand"

// Gender defines the person's gender
type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
	Other  Gender = "Other"
)

func GetRandomGender() Gender {
	randomGender := rand.Intn(100)
	gender := Other
	switch {
	case randomGender < 49:
		gender = Male
	case randomGender < 98:
		gender = Female
	}
	return gender
}
