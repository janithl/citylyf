package people

import (
	"math/rand"
)

// RelationshipStatus defines if the person is married etc.
type RelationshipStatus string

const (
	Single   RelationshipStatus = "Single"
	Married  RelationshipStatus = "Married"
	Divorced RelationshipStatus = "Divorced"
	Widowed  RelationshipStatus = "Widowed"
)

func getRelationshipStatus(age int) RelationshipStatus {
	if age < ageOfAdulthood {
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
