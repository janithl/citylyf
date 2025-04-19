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
	Scale               float64
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

	if b.Scale > 1 {
		image := ebiten.NewImage(b.Width, b.Height)
		vector.DrawFilledRect(image, 0, 0, float32(b.Width), float32(b.Height), btnColor, false)
		ebitenutil.DebugPrintAt(image, b.Label, 10, 4)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(b.Scale, b.Scale)
		op.GeoM.Translate(float64(b.X), float64(b.Y))
		screen.DrawImage(image, op)

	} else {
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), btnColor, false)
		ebitenutil.DebugPrintAt(screen, b.Label, b.X+10, b.Y+4)
	}
}

func (b *Button) Update() {
	mouseX, mouseY := ebiten.CursorPosition()
	isPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Check if mouse is over the button
	if b.Scale > 1 {
		height := int(float64(b.Height) * b.Scale)
		width := int(float64(b.Width) * b.Scale)
		b.isHovered = mouseX > b.X && mouseX < b.X+width && mouseY > b.Y && mouseY < b.Y+height
	} else {
		b.isHovered = mouseX > b.X && mouseX < b.X+b.Width && mouseY > b.Y && mouseY < b.Y+b.Height
	}

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
