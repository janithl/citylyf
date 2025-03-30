package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type Tooltip struct {
	X, Y, Height, Width, Padding, Margin int
	visible                              bool
	Text                                 string
}

func (t *Tooltip) Draw(screen *ebiten.Image) {
	if !t.visible {
		return
	}

	x, y := t.X+t.Margin, t.Y+t.Margin
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(t.Width), float32(t.Height), colour.DarkSemiBlack, false)
	ebitenutil.DebugPrintAt(screen, t.Text, x+t.Padding, y+t.Padding)
}

func (t *Tooltip) Update(visible bool) {
	t.visible = visible
}
