package control

import "github.com/hajimehoshi/ebiten/v2"

type GraphWindow struct {
	dataSource   func() []float64 // Function to dynamically fetch data
	Window       *Window
	frameCounter int
}

func (gw *GraphWindow) Update() error {
	gw.frameCounter++
	if gw.frameCounter >= 60 { // update every second
		gw.frameCounter = 0

		// Find and update the existing graph
		if graph, ok := gw.Window.Children[0].(*Graph); ok {
			graph.Data = gw.dataSource() // Get fresh data from source
		}
	}
	gw.Window.Update()
	return nil
}

func (gw *GraphWindow) Draw(screen *ebiten.Image) {
	gw.Window.Draw(screen)
}

// NewGraphWindow creates a new graph window instance
func NewGraphWindow(x, y, width, height int, title string, closeFunc func(string), dataSource func() []float64) *GraphWindow {
	window := NewWindow(x, y, width, height, title, closeFunc)
	window.AddChild(&Graph{
		X:      0,
		Y:      0,
		Width:  float32(width),
		Height: float32(height - 24),
		Data:   dataSource(),
	})
	return &GraphWindow{
		Window:     window,
		dataSource: dataSource,
	}
}
