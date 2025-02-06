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
	OnClick             func(index int)
}

func (tl *TextList) Update() {
	for i, btn := range tl.Buttons {
		if i < len(tl.Items) {
			btn.Label = tl.Items[i]
		}
		btn.Update()
	}
}

func (tl *TextList) Draw(screen *ebiten.Image) {
	for _, btn := range tl.Buttons {
		btn.Draw(screen)
	}
}

func (tl *TextList) SetOffset(x, y int) {
	tl.X = x
	tl.Y = y
	for i, btn := range tl.Buttons {
		btn.SetOffset(x, y+i*buttonHeight)
	}
}

// createButtons initializes buttons for the current page
func (tl *TextList) createButtons() {
	tl.Buttons = nil // Reset buttons

	for i := 0; i < tl.Height/buttonHeight; i++ {
		btn := &Button{
			X: tl.X, Y: tl.Y + i*buttonHeight, Width: tl.Width, Height: buttonHeight,
			Label:      "",
			Color:      colour.Transparent,
			HoverColor: colour.SemiBlack,
		}
		btn.OnClick = func() { tl.OnClick(i) }
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
