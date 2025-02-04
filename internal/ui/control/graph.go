package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type Graph struct {
	Title               string // Graph title
	X, Y, Width, Height float32
	Data                []float64 // Y-values of the graph
	MaxValue            float64   // Max Y-value for scaling
}

// Draws the graph on the screen
func (g *Graph) Draw(screen *ebiten.Image) {
	// Draw background and horizontal grid lines
	vector.DrawFilledRect(screen, (g.X), (g.Y), (g.Width), (g.Height), colour.Black, true)
	vector.StrokeLine(screen, (g.X), (g.Y + g.Height*0.5), (g.X)+(g.Width), (g.Y + g.Height*0.5), 1.0, colour.White, true)
	vector.StrokeLine(screen, (g.X), (g.Y + g.Height*0.25), (g.X)+(g.Width), (g.Y + g.Height*0.25), 1.0, colour.Gray, true)
	vector.StrokeLine(screen, (g.X), (g.Y + g.Height*0.75), (g.X)+(g.Width), (g.Y + g.Height*0.75), 1.0, colour.Gray, true)

	if len(g.Data) < 2 {
		return // Not enough points to draw a line
	}

	// Draw data lines
	pointCount := len(g.Data)
	step := g.Width / float32(pointCount-1) // Space between points

	for i := 0; i < pointCount-1.0; i++ {
		x1 := (g.X) + step*float32(i)
		y1 := (g.Y + g.Height) - float32(g.Data[i]/(g.MaxValue+10))*(g.Height) // Scale Y
		x2 := (g.X) + step*float32(i+1)
		y2 := (g.Y + g.Height) - float32(g.Data[i+1]/(g.MaxValue+10))*(g.Height)

		vector.StrokeLine(screen, x1, y1, x2, y2, 2.0, colour.Green, true)
	}

	ebitenutil.DebugPrintAt(screen, g.Title, int(g.X)+4, int(g.Y))
}
