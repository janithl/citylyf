package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type MapWindow struct {
	x, y, width, height int
	stepperLabels       []*Label
	steppers            []*Stepper
	buttons             []*Button
}

func (mw *MapWindow) Update() {
	for i := range mw.stepperLabels {
		mw.stepperLabels[i].Update()
	}
	for i := range mw.steppers {
		mw.steppers[i].Update()
	}
	for i := range mw.buttons {
		mw.buttons[i].Update()
	}
}

func (mw *MapWindow) Draw(screen *ebiten.Image) {
	for i := range mw.stepperLabels {
		mw.stepperLabels[i].Draw(screen)
	}
	for i := range mw.steppers {
		mw.steppers[i].Draw(screen)
	}
	for i := range mw.buttons {
		mw.buttons[i].Draw(screen)
	}
}

func (mw *MapWindow) SetOffset(x, y int) {
	mw.x = x
	mw.y = y
	for i, stepperLabel := range mw.stepperLabels {
		stepperLabel.SetOffset(x+4, y+4+(i*buttonHeight))
	}
	for i, stepper := range mw.steppers {
		stepper.SetOffset(x+mw.width/2, y+(i*buttonHeight))
	}
	for i, btn := range mw.buttons {
		btn.SetOffset(x, y+(mw.height-i*buttonHeight))
	}
}

func NewMapWindow(x, y, width, height int, closeFunc func()) *MapWindow {
	steppersLabels := []*Label{
		{X: 0, Y: 0, Text: "Mountain Peaks"},
		{X: 0, Y: 0, Text: "Mountain Ranges"},
		{X: 0, Y: 0, Text: "Cliffs"},
	}
	steppers := []*Stepper{
		NewStepper(0, 0, 50, 90, PercentageStepper, func(i int) {}),
		NewStepper(0, 0, 50, 90, PercentageStepper, func(i int) {}),
		NewStepper(0, 0, 50, 90, PercentageStepper, func(i int) {}),
	}
	buttons := []*Button{
		{Label: "Done", X: 0, Y: 0, Width: 240, Height: buttonHeight, Color: colour.Black, HoverColor: colour.Red, OnClick: closeFunc},
		{Label: "Regenerate Map", X: 0, Y: 0, Width: 240, Height: buttonHeight, Color: colour.Black, HoverColor: colour.DarkGreen, OnClick: entities.Sim.RegenerateMap},
	}

	return &MapWindow{
		x:             x,
		y:             y,
		width:         width,
		height:        height,
		stepperLabels: steppersLabels,
		steppers:      steppers,
		buttons:       buttons,
	}
}
