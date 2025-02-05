package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
)

const (
	titleBarHeight   = 24
	closeButtonWidth = 36
)

// Window represents a movable, dismissible UI window
type Window struct {
	X, Y, Width, Height      int
	Title                    string
	IsDragging               bool
	IsVisible                bool
	DragOffsetX, DragOffsetY int
	Children                 []Renderable // Child controls
	closeButton              *Button
}

// NewWindow creates a new window instance
func NewWindow(x, y, width, height int, title string, closeFunc func(string)) *Window {
	return &Window{
		X:         x,
		Y:         y,
		Width:     width,
		Height:    height,
		Title:     title,
		IsVisible: true,
		closeButton: &Button{
			Label:      " X ",
			X:          x,
			Y:          y,
			Width:      closeButtonWidth,
			Height:     titleBarHeight,
			Color:      colour.Black,
			HoverColor: colour.Red,
			OnClick:    func() { closeFunc(title) },
		},
	}
}

// Update handles window dragging and dismissal
func (w *Window) Update() {
	if !w.IsVisible {
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()
	isPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Handle dragging
	if isPressed {
		if !w.IsDragging && mouseX > w.X && mouseX < w.X+w.Width && mouseY > w.Y && mouseY < w.Y+30 {
			w.IsDragging = true
			w.DragOffsetX = mouseX - w.X
			w.DragOffsetY = mouseY - w.Y
		}
	} else {
		w.IsDragging = false
	}

	if w.IsDragging {
		w.X = mouseX - w.DragOffsetX
		w.Y = mouseY - w.DragOffsetY
	}

	// update window position to its contents
	w.closeButton.X, w.closeButton.Y = w.X, w.Y
	w.closeButton.Update()

	// update child controls
	for _, child := range w.Children {
		child.SetOffset(w.X, w.Y+titleBarHeight) // Offset for title bar
		child.Update()
	}
}

// Draw renders the window and child controls
func (w *Window) Draw(screen *ebiten.Image) {
	if !w.IsVisible {
		return
	}

	// Draw window background
	vector.DrawFilledRect(screen, float32(w.X), float32(w.Y), float32(w.Width), float32(w.Height), colour.DarkGray, true)

	// Draw title bar
	vector.DrawFilledRect(screen, float32(w.X), float32(w.Y), float32(w.Width), titleBarHeight, colour.Black, true)
	ebitenutil.DebugPrintAt(screen, w.Title, w.X+closeButtonWidth+10, w.Y+4)

	// Draw close button
	w.closeButton.Draw(screen)

	// Draw child controls
	for _, child := range w.Children {
		child.Draw(screen)
	}
}

// AddChild adds a child control to the window
func (w *Window) AddChild(child Renderable) {
	w.Children = append(w.Children, child)
}

// ClearChildren clears all the child controls
func (w *Window) ClearChildren() {
	w.Children = []Renderable{}
}

// CloseWindow closes the window by making IsVisible flase
func (w *Window) CloseWindow() {
	w.IsVisible = false
}
