package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/ui/colour"
)

const (
	buttonHeight = 24
)

type TextList struct {
	X, Y, Width, Height int
	Items               []string
	Buttons             []*Button
}

func (tl *TextList) Update() {
	for _, btn := range tl.Buttons {
		btn.SetOffset(tl.X, tl.Y)
		btn.Update()
	}
}

func (tl *TextList) Draw(screen *ebiten.Image) {
	for _, btn := range tl.Buttons {
		btn.Draw(screen)
	}
}

func (tl *TextList) SetOffset(x, y int) {
	tl.X += x
	tl.Y += y
}

// createButtons initializes buttons for the current page
func (tl *TextList) createButtons() {
	tl.Buttons = nil // Reset buttons

	for i, text := range tl.Items {
		btn := &Button{
			X: tl.X, Y: tl.Y + i*buttonHeight, Width: tl.Width, Height: buttonHeight,
			Label:      text,
			Color:      colour.Transparent,
			HoverColor: colour.SemiBlack,
			OnClick:    func() { println("Clicked:", text) },
		}
		tl.Buttons = append(tl.Buttons, btn)
	}
}

// NewTextList creates a list
func NewTextList(x, y, width, height int, items []string) *TextList {
	tl := &TextList{
		X: x, Y: y, Width: width, Height: height,
		Items: items,
	}

	tl.createButtons()
	return tl
}
