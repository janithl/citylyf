package entities

// CareerLevel defines the levels in a person's career
type CareerLevel string

const (
	Unemployed     CareerLevel = "Unemployed"
	EntryLevel     CareerLevel = "Entry Level"
	MidLevel       CareerLevel = "Mid Level"
	SeniorLevel    CareerLevel = "Senior Level"
	ExecutiveLevel CareerLevel = "Executive Level"
)

var CareerLevels = []CareerLevel{Unemployed, EntryLevel, MidLevel, SeniorLevel, ExecutiveLevel}
