package control

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type Stepper struct {
	X, Y           int
	currentNumber  int
	maxNumber      int
	StepperType    StepperType
	decreaseButton *Button
	increaseButton *Button
	onChange       func(newValue int) // Callback when value changes
}

// Update processes button clicks
func (s *Stepper) Update() {
	s.decreaseButton.Update()
	s.increaseButton.Update()
}

// Draw renders the stepper
func (s *Stepper) Draw(screen *ebiten.Image) {
	s.decreaseButton.Draw(screen)
	s.increaseButton.Draw(screen)
	var text string
	if s.StepperType == NumberStepper {
		text = fmt.Sprintf("%02d/%02d", s.currentNumber, s.maxNumber)
	} else if s.StepperType == PercentageStepper {
		text = fmt.Sprintf(" %02d%% ", s.currentNumber)
	}
	ebitenutil.DebugPrintAt(screen, text, s.X+buttonWidth+2, s.Y+4)
}

// SetOffset moves the stepper when the parent window moves
func (s *Stepper) SetOffset(x, y int) {
	s.X = x
	s.Y = y
	s.decreaseButton.SetOffset(x, y)
	s.increaseButton.SetOffset(x+2*buttonWidth, y)
}

// SetMaxNumber sets the max number
func (s *Stepper) SetMaxNumber(maxNumber int) {
	if maxNumber > 0 {
		s.maxNumber = maxNumber
	}
}

// SetCurrentNumber sets the current stepper value
func (s *Stepper) SetCurrentNumber(currentNumber int) {
	s.currentNumber = currentNumber
}

// NewStepper creates a new stepper with an optional callback
func NewStepper(x, y, currentNumber, maxNumber int, stepperType StepperType, onChange func(int)) *Stepper {
	if maxNumber < currentNumber {
		maxNumber = currentNumber // Ensure valid max
	}

	stepper := &Stepper{
		X: x, Y: y,
		StepperType:   stepperType,
		currentNumber: currentNumber,
		maxNumber:     maxNumber,
		onChange:      onChange,
	}

	leftLabel, rightLabel := " < ", " > "
	leftIncrement, rightIncrement := -1, 1
	if stepperType == PercentageStepper {
		leftLabel, rightLabel = " - ", " + "
		leftIncrement, rightIncrement = -10, 10
	}

	stepper.decreaseButton = &Button{
		X: x, Y: y, Width: buttonWidth, Height: buttonHeight,
		Label: leftLabel, Color: colour.Transparent, HoverColor: colour.SemiBlack,
		OnClick: func() {
			if stepper.currentNumber > 1 {
				stepper.currentNumber += leftIncrement
				stepper.onChange(stepper.currentNumber)
			}
		},
	}

	stepper.increaseButton = &Button{
		X: x + 2*buttonWidth, Y: y, Width: buttonWidth, Height: buttonHeight,
		Label: rightLabel, Color: colour.Transparent, HoverColor: colour.SemiBlack,
		OnClick: func() {
			if stepper.currentNumber < stepper.maxNumber {
				stepper.currentNumber += rightIncrement
				stepper.onChange(stepper.currentNumber)
			}
		},
	}

	return stepper
}
