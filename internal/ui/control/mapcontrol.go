package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type MapControl struct {
	x, y, width, height            int
	peakPerc, rangePerc, cliffPerc int
	layoutGrid                     *Grid
}

func (mc *MapControl) Update() {
	mc.layoutGrid.Update()
}

func (mc *MapControl) Draw(screen *ebiten.Image) {
	mc.layoutGrid.Draw(screen)
}

func (mc *MapControl) SetOffset(x, y int) {
	mc.x = x
	mc.y = y
	mc.layoutGrid.SetOffset(x, y)
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
	mc.layoutGrid.Children[2][3].(*Stepper).SetCurrentNumber(mc.peakPerc)
	mc.layoutGrid.Children[3][3].(*Stepper).SetCurrentNumber(mc.rangePerc)
	mc.layoutGrid.Children[4][3].(*Stepper).SetCurrentNumber(mc.cliffPerc)
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

func (mc *MapControl) saveCityName() {
	if mc.layoutGrid.Children[0][2].(*TextInput).Text != "" {
		entities.Sim.CityName = mc.layoutGrid.Children[0][2].(*TextInput).Text
	} else {
		entities.Sim.CityName = "UnnamedCity"
	}
}

func NewMapControl(x, y, width, height int, closeFunc func()) *MapControl {
	mc := &MapControl{
		x:          x,
		y:          y,
		width:      width,
		height:     height,
		peakPerc:   30,
		rangePerc:  10,
		cliffPerc:  10,
		layoutGrid: NewGrid(x, y, width, height, 6, 8),
	}

	mc.layoutGrid.Children[0][0] = &Label{X: 0, Y: 0, Padding: 4, Text: "City Name"}
	mc.layoutGrid.Children[0][2] = &TextInput{X: 0, Y: 0, Padding: 4, Width: 2 * width / 3, Height: buttonHeight, Focused: true}

	mc.layoutGrid.Children[2][0] = &Label{X: 0, Y: 0, Padding: 4, Text: "Mountain Peaks"}
	mc.layoutGrid.Children[2][3] = NewStepper(0, 0, mc.peakPerc, 90, PercentageStepper, func(i int) { mc.setPerc("peak", i) })

	mc.layoutGrid.Children[3][0] = &Label{X: 0, Y: 0, Padding: 4, Text: "Mountain Ranges"}
	mc.layoutGrid.Children[3][3] = NewStepper(0, 0, mc.rangePerc, 90, PercentageStepper, func(i int) { mc.setPerc("range", i) })

	mc.layoutGrid.Children[4][0] = &Label{X: 0, Y: 0, Padding: 4, Text: "Cliffs"}
	mc.layoutGrid.Children[4][3] = NewStepper(0, 0, mc.cliffPerc, 90, PercentageStepper, func(i int) { mc.setPerc("cliff", i) })

	mc.layoutGrid.Children[6][0] = &Button{Label: "  Regen", X: 0, Y: 0, Width: width / 3, Height: buttonHeight, Color: colour.Transparent, HoverColor: colour.DarkCyan, OnClick: mc.regenerateMap}
	mc.layoutGrid.Children[6][2] = &Button{Label: "  Reset", X: 0, Y: 0, Width: width / 3, Height: buttonHeight, Color: colour.Transparent, HoverColor: colour.DarkMagenta, OnClick: mc.resetValues}
	mc.layoutGrid.Children[6][4] = &Button{Label: "  Done", X: 0, Y: 0, Width: width / 3, Height: buttonHeight, Color: colour.Transparent, HoverColor: colour.Red, OnClick: func() { mc.saveCityName(); closeFunc() }}

	return mc
}
