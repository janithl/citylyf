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
	text := fmt.Sprintf("%02d/%02d", s.currentNumber, s.maxNumber)
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

// NewStepper creates a new stepper with an optional callback
func NewStepper(x, y, maxNumber int, onChange func(int)) *Stepper {
	if maxNumber < 1 {
		maxNumber = 1 // Ensure valid max
	}

	stepper := &Stepper{
		X: x, Y: y,
		currentNumber: 1, maxNumber: maxNumber,
		onChange: onChange,
	}

	stepper.decreaseButton = &Button{
		X: x, Y: y, Width: buttonWidth, Height: buttonHeight,
		Label: " < ", Color: colour.Transparent, HoverColor: colour.SemiBlack,
		OnClick: func() {
			if stepper.currentNumber > 1 {
				stepper.currentNumber--
				stepper.onChange(stepper.currentNumber)
			}
		},
	}

	stepper.increaseButton = &Button{
		X: x + 2*buttonWidth, Y: y, Width: buttonWidth, Height: buttonHeight,
		Label: " > ", Color: colour.Transparent, HoverColor: colour.SemiBlack,
		OnClick: func() {
			if stepper.currentNumber < stepper.maxNumber {
				stepper.currentNumber++
				stepper.onChange(stepper.currentNumber)
			}
		},
	}

	return stepper
}
