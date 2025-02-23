package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Label struct {
	X, Y, Padding int
	Text          string
}

func (l *Label) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, l.Text, l.X+l.Padding, l.Y+l.Padding)
}

func (l *Label) Update() {}

func (l *Label) SetOffset(x, y int) {
	l.X = x
	l.Y = y
}
