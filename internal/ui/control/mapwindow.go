package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type MapWindow struct {
	x, y, width, height            int
	peakPerc, rangePerc, cliffPerc int
	stepperLabels                  []*Label
	steppers                       []*Stepper
	buttons                        []*Button
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

func (mw *MapWindow) setPerc(perc string, value int) {
	switch perc {
	case "peak":
		mw.peakPerc = value
	case "range":
		mw.rangePerc = value
	case "cliff":
		mw.cliffPerc = value
	}
}

func (mw *MapWindow) regenerateMap() {
	peakProb := 0.001 + 0.00001*float64(mw.peakPerc)
	rangeProb := 0.0 + 0.0001*float64(mw.rangePerc)
	cliffProb := 0.0 + 0.0002*float64(mw.cliffPerc)
	entities.Sim.RegenerateMap(peakProb, rangeProb, cliffProb)
}

func NewMapWindow(x, y, width, height int, closeFunc func()) *MapWindow {
	mw := &MapWindow{
		x:         x,
		y:         y,
		width:     width,
		height:    height,
		peakPerc:  50,
		rangePerc: 50,
		cliffPerc: 50,
	}

	mw.stepperLabels = []*Label{
		{X: 0, Y: 0, Text: "Mountain Peaks"},
		{X: 0, Y: 0, Text: "Mountain Ranges"},
		{X: 0, Y: 0, Text: "Cliffs"},
	}

	mw.steppers = []*Stepper{
		NewStepper(0, 0, mw.peakPerc, 90, PercentageStepper, func(i int) { mw.setPerc("peak", i) }),
		NewStepper(0, 0, mw.rangePerc, 90, PercentageStepper, func(i int) { mw.setPerc("range", i) }),
		NewStepper(0, 0, mw.cliffPerc, 90, PercentageStepper, func(i int) { mw.setPerc("cliff", i) }),
	}

	mw.buttons = []*Button{
		{Label: "Done", X: 0, Y: 0, Width: 240, Height: buttonHeight, Color: colour.Black, HoverColor: colour.Red, OnClick: closeFunc},
		{Label: "Regenerate Map", X: 0, Y: 0, Width: 240, Height: buttonHeight, Color: colour.Black, HoverColor: colour.DarkGreen, OnClick: mw.regenerateMap},
	}

	return mw
}
