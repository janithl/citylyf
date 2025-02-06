package control

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type Graph struct {
	X, Y, Width, Height float32
	Data                []float64 // Y-values of the graph
}

// Draws the graph on the screen
func (g *Graph) Draw(screen *ebiten.Image) {
	// Draw horizontal grid lines
	for i := float32(0.0); i <= 1.0; i += 0.25 {
		vector.StrokeLine(screen, g.X, g.Y+(g.Height*i), g.X+g.Width, g.Y+(g.Height*i), 1.0, colour.Black, true)
	}

	if len(g.Data) < 2 {
		return // Not enough points to draw a line
	}

	// Determine the min and max values in the data set
	minValue, maxValue, lastValue := 0.0, 0.0, 0.0
	for _, val := range g.Data {
		if val > maxValue {
			maxValue = val
		}
		if val < minValue {
			minValue = val
		}
		lastValue = val
	}
	valueRange := maxValue - minValue

	// Calculate step size (spacing between points)
	pointCount := len(g.Data)
	step := g.Width / float32(math.Max(float64(pointCount-1), 8))

	// Draw data lines
	for i := 0; i < pointCount-1; i++ {
		x1 := g.X + step*float32(i)
		x2 := g.X + step*float32(i+1)

		// Scale Y values to fit the graph
		y1 := g.Y + g.Height - float32((g.Data[i]-minValue)/valueRange)*g.Height
		y2 := g.Y + g.Height - float32((g.Data[i+1]-minValue)/valueRange)*g.Height

		vector.StrokeLine(screen, x1, y1, x2, y2, 2.0, colour.Green, true)
	}

	// add value label
	label := fmt.Sprintf("%.2f", lastValue)
	ebitenutil.DebugPrintAt(screen, label, int(g.X)+8, int(g.Y+g.Height)-12)
}

func (g *Graph) Update() {}

func (g *Graph) SetOffset(x, y int) {
	g.X = float32(x)
	g.Y = float32(y)
}
