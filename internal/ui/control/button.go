package control

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	Label               string
	X, Y, Width, Height int
	Color               color.RGBA
	OnClick             func() // Callback function

	HoverColor color.RGBA
	isHovered  bool
	wasPressed bool // to prevent multiple clicks
}

func (b *Button) Draw(screen *ebiten.Image) {
	btnColor := b.Color
	if b.isHovered {
		btnColor = b.HoverColor
	}

	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), btnColor, true)
	ebitenutil.DebugPrintAt(screen, b.Label, b.X+10, b.Y+4)
}

func (b *Button) Update() {
	mouseX, mouseY := ebiten.CursorPosition()
	isPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Check if mouse is over the button
	b.isHovered = mouseX > b.X && mouseX < b.X+b.Width && mouseY > b.Y && mouseY < b.Y+b.Height

	// Only trigger OnClick when mouse is released after a press
	if b.isHovered && !isPressed && b.wasPressed {
		b.OnClick()
	}

	// Track previous mouse state
	b.wasPressed = isPressed
}

func (b *Button) SetOffset(x, y int) {
	b.X = x
	b.Y = y
}
