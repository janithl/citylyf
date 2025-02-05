package control

import "github.com/hajimehoshi/ebiten/v2"

type GraphWindow struct {
	Data   *[]float64
	Window *Window
}

func (g *GraphWindow) Update() error {
	g.Window.ClearChildren()
	g.Window.AddChild(&Graph{
		X:      0,
		Y:      2,
		Width:  float32(g.Window.Width),
		Height: float32(g.Window.Height - 30),
		Data:   *g.Data,
	})
	g.Window.Update()
	return nil
}

func (g *GraphWindow) Draw(screen *ebiten.Image) {
	g.Window.Draw(screen)
}
