package control

import "github.com/hajimehoshi/ebiten/v2"

type GraphWindow struct {
	Data   *[]float64
	Window *Window
}

func (g *GraphWindow) Update() error {
	// Find and update the existing graph
	if graph, ok := g.Window.Children[0].(*Graph); ok {
		graph.Data = *g.Data
	}
	g.Window.Update()
	return nil
}

func (g *GraphWindow) Draw(screen *ebiten.Image) {
	g.Window.Draw(screen)
}

// NewGraphWindow creates a new graph window instance
func NewGraphWindow(x, y, width, height int, title string, closeFunc func(string), data *[]float64) *GraphWindow {
	window := NewWindow(x, y, width, height, title, closeFunc)
	window.AddChild(&Graph{
		X:      0,
		Y:      0,
		Width:  float32(width),
		Height: float32(height - 24),
		Data:   *data,
	})
	return &GraphWindow{
		Window: window,
		Data:   data,
	}
}
