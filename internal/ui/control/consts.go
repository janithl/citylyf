package control

const (
	buttonHeight    = 24
	buttonWidth     = 36
	titleBarHeight  = 24
	menuEntryHeight = 72
)

type GraphType int

const (
	Int        GraphType = 0
	Float      GraphType = 1
	Percentage GraphType = 2
	Currency   GraphType = 3
)

type StepperType int

const (
	NumberStepper     StepperType = 0
	PercentageStepper StepperType = 1
)
