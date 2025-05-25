package control

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type TextInput struct {
	X, Y, Padding int
	Width, Height int
	Text          string
	Focused       bool
}

func (ti *TextInput) Draw(screen *ebiten.Image) {
	cursor := ""
	if ti.Focused {
		cursor = "_"
	}

	vector.StrokeRect(screen, float32(ti.X), float32(ti.Y), float32(ti.Width), float32(ti.Height), 1, colour.White, false)
	ebitenutil.DebugPrintAt(screen, ti.Text+cursor, ti.X+ti.Padding, ti.Y+ti.Padding)

}

func (ti *TextInput) Update() {
	mx, my := ebiten.CursorPosition()
	clicked := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Detect click within bounds to toggle focus
	if clicked && mx >= ti.X && mx <= ti.X+ti.Width && my >= ti.Y && my <= ti.Y+ti.Height {
		ti.Focused = true
	} else if clicked {
		ti.Focused = false
	}

	if ti.Focused {
		var inputChars []rune
		inputChars = ebiten.AppendInputChars(inputChars[:0])
		for _, char := range inputChars {
			if char >= 32 && char <= 126 {
				ti.Text += string(char)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(ti.Text) > 0 {
			ti.Text = strings.TrimSuffix(ti.Text, string(ti.Text[len(ti.Text)-1]))
		}
	}
}

func (ti *TextInput) SetOffset(x, y int) {
	ti.X = x
	ti.Y = y
}
