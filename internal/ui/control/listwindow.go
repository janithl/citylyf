package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
)

type ListWindow struct {
	dataSource   func() []Statable // Function to dynamically fetch data
	Window       *Window
	frameCounter int
}

func (lw *ListWindow) Update() error {
	lw.frameCounter++
	if lw.frameCounter >= 60 { // update every second
		lw.frameCounter = 0

		// Find and update the existing text list
		if list, ok := lw.Window.Children[0].(*TextList); ok {
			entities.Sim.Mutex.RLock()
			list.UpdateItems(lw.dataSource()) // Updates text without resetting buttons
			entities.Sim.Mutex.RUnlock()
		}
	}
	lw.Window.Update()
	return nil
}

func (lw *ListWindow) Draw(screen *ebiten.Image) {
	lw.Window.Draw(screen)
}

// NewListWindow creates a new graph window instance
func NewListWindow(x, y, width, height int, title string, closeFunc func(string), clickFunc func(string, int), dataSource func() []Statable) *ListWindow {
	window := NewWindow(x, y, width, height, title, closeFunc)
	textlist := NewTextList(0, 0, width, height-titleBarHeight, dataSource())
	textlist.OnClick = func(index int) {
		clickFunc(title, index)
	}
	window.AddChild(textlist)
	return &ListWindow{
		Window:     window,
		dataSource: dataSource,
	}
}
