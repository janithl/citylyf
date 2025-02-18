package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type MapWindow struct {
	Window *Window
}

func (mw *MapWindow) Update() error {
	mw.Window.Update()
	return nil
}

func (mw *MapWindow) Draw(screen *ebiten.Image) {
	mw.Window.Draw(screen)

}

func NewMapWindow(x, y, width, height int, closeFunc func(string)) *MapWindow {
	window := NewWindow(x, y, width, height, "Map Control", closeFunc)
	window.AddChild(&Button{Label: "Regenerate Map", X: 0, Y: 0, Width: 200, Height: buttonHeight, Color: colour.Black, HoverColor: colour.DarkGreen, OnClick: entities.Sim.RegenerateMap})

	return &MapWindow{
		Window: window,
	}
}
