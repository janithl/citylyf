package entities

import "math/rand"

// Industry defines the industry of the business
type Industry string

const (
	Agriculture        Industry = "Agriculture"
	Automobile         Industry = "Automobile"
	Construction       Industry = "Construction"
	Creative           Industry = "Creative"
	Education          Industry = "Education"
	Energy             Industry = "Energy"
	Finance            Industry = "Finance"
	Healthcare         Industry = "Healthcare"
	Retail             Industry = "Retail"
	Technology         Industry = "Technology"
	Telecommunications Industry = "Telecommunications"
)

var industries = []Industry{
	Agriculture, Automobile, Construction, Education, Energy, Finance, Healthcare, Retail, Technology, Telecommunications,
}

func GetRandomIndustry() Industry {
	return industries[rand.Intn(len(industries))]
}
