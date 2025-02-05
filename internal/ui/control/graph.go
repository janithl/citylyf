package control

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type Graph struct {
	Title               string // Graph title
	X, Y, Width, Height float32
	Data                []float64 // Y-values of the graph
}

// Draws the graph on the screen
func (g *Graph) Draw(screen *ebiten.Image) {
	// Draw horizontal grid lines
	vector.StrokeLine(screen, g.X, g.Y, g.X+g.Width, g.Y, 1.0, colour.Black, true)
	vector.StrokeLine(screen, g.X, g.Y+g.Height*0.25, g.X+g.Width, g.Y+g.Height*0.25, 1.0, colour.Black, true)
	vector.StrokeLine(screen, g.X, g.Y+g.Height*0.5, g.X+g.Width, g.Y+g.Height*0.5, 2.0, colour.Black, true)
	vector.StrokeLine(screen, g.X, g.Y+g.Height*0.75, g.X+g.Width, g.Y+g.Height*0.75, 1.0, colour.Black, true)
	vector.StrokeLine(screen, g.X, g.Y+g.Height, g.X+g.Width, g.Y+g.Height, 1.0, colour.Black, true)

	if len(g.Data) < 2 {
		return // Not enough points to draw a line
	}

	// Draw data lines
	pointCount := len(g.Data)
	step := g.Width / float32(math.Max(float64(pointCount-1), 8)) // Space between points
	maxValue := 0.0
	for _, val := range g.Data {
		if val > maxValue {
			maxValue = val
		}
	}

	for i := 0; i < pointCount-1.0; i++ {
		x1 := g.X + step*float32(i)
		y1 := g.Y + g.Height - float32(g.Data[i]/maxValue)*g.Height // Scale Y
		x2 := g.X + step*float32(i+1)
		y2 := g.Y + g.Height - float32(g.Data[i+1]/maxValue)*g.Height

		vector.StrokeLine(screen, x1, y1, x2, y2, 2.0, colour.Green, true)
	}

	ebitenutil.DebugPrintAt(screen, g.Title, int(g.X)+4, int(g.Y))
}

func (g *Graph) Update() {}

func (g *Graph) SetOffset(x, y int) {
	g.X += float32(x)
	g.Y += float32(y)
}
