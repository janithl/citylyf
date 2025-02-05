package control

import "github.com/hajimehoshi/ebiten/v2"

// Renderable interface allows child elements to be rendered inside the window
type Renderable interface {
	SetOffset(x, y int)
	Draw(screen *ebiten.Image)
	Update()
}
