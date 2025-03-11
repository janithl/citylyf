package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type MapControl struct {
	x, y, width, height            int
	peakPerc, rangePerc, cliffPerc int
	stepperLabels                  []*Label
	steppers                       []*Stepper
	buttons                        []*Button
}

func (mc *MapControl) Update() {
	for i := range mc.stepperLabels {
		mc.stepperLabels[i].Update()
	}
	for i := range mc.steppers {
		mc.steppers[i].Update()
	}
	for i := range mc.buttons {
		mc.buttons[i].Update()
	}
}

func (mc *MapControl) Draw(screen *ebiten.Image) {
	for i := range mc.stepperLabels {
		mc.stepperLabels[i].Draw(screen)
	}
	for i := range mc.steppers {
		mc.steppers[i].Draw(screen)
	}
	for i := range mc.buttons {
		mc.buttons[i].Draw(screen)
	}
}

func (mc *MapControl) SetOffset(x, y int) {
	mc.x = x
	mc.y = y
	for i, stepperLabel := range mc.stepperLabels {
		stepperLabel.SetOffset(x, y+(i*buttonHeight))
	}
	for i, stepper := range mc.steppers {
		stepper.SetOffset(x+mc.width/2, y+(i*buttonHeight))
	}
	for i, btn := range mc.buttons {
		btn.SetOffset(x+i*(mc.width/3), y+(4*buttonHeight))
	}
}

func (mc *MapControl) setPerc(perc string, value int) {
	switch perc {
	case "peak":
		mc.peakPerc = value
	case "range":
		mc.rangePerc = value
	case "cliff":
		mc.cliffPerc = value
	}
	mc.regenerateMap()
}

func (mc *MapControl) resetValues() {
	mc.peakPerc = 30
	mc.rangePerc = 10
	mc.cliffPerc = 10
	mc.steppers[0].SetCurrentNumber(mc.peakPerc)
	mc.steppers[1].SetCurrentNumber(mc.rangePerc)
	mc.steppers[2].SetCurrentNumber(mc.cliffPerc)
	mc.regenerateMap()
}

func (mc *MapControl) regenerateMap() {
	peakProb := 0.00005 * float64(mc.peakPerc)
	rangeProb := 0.0005 * float64(mc.rangePerc)
	cliffProb := 0.01 * float64(mc.cliffPerc)
	entities.Sim.Mutex.Lock()
	entities.Sim.RegenerateMap(peakProb, rangeProb, cliffProb)
	entities.Sim.Mutex.Unlock()
}

func NewMapControl(x, y, width, height int, closeFunc func()) *MapControl {
	mc := &MapControl{
		x:         x,
		y:         y,
		width:     width,
		height:    height,
		peakPerc:  30,
		rangePerc: 10,
		cliffPerc: 10,
	}

	mc.stepperLabels = []*Label{
		{X: 0, Y: 0, Padding: 4, Text: "Mountain Peaks"},
		{X: 0, Y: 0, Padding: 4, Text: "Mountain Ranges"},
		{X: 0, Y: 0, Padding: 4, Text: "Cliffs"},
	}

	mc.steppers = []*Stepper{
		NewStepper(0, 0, mc.peakPerc, 90, PercentageStepper, func(i int) { mc.setPerc("peak", i) }),
		NewStepper(0, 0, mc.rangePerc, 90, PercentageStepper, func(i int) { mc.setPerc("range", i) }),
		NewStepper(0, 0, mc.cliffPerc, 90, PercentageStepper, func(i int) { mc.setPerc("cliff", i) }),
	}

	mc.buttons = []*Button{
		{Label: "  Regen", X: 0, Y: 0, Width: width / 3, Height: buttonHeight, Color: colour.Transparent, HoverColor: colour.DarkCyan, OnClick: mc.regenerateMap},
		{Label: "  Reset", X: 0, Y: 0, Width: width / 3, Height: buttonHeight, Color: colour.Transparent, HoverColor: colour.DarkMagenta, OnClick: mc.resetValues},
		{Label: "  Done", X: 0, Y: 0, Width: width / 3, Height: buttonHeight, Color: colour.Transparent, HoverColor: colour.Red, OnClick: closeFunc},
	}

	return mc
}
