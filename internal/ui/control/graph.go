package control

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
	"github.com/janithl/citylyf/internal/utils"
)

type GraphType int

const (
	Int        GraphType = 0
	Float      GraphType = 1
	Percentage GraphType = 2
	Currency   GraphType = 3
)

type Graph struct {
	x, y, width, height float32
	graphType           GraphType
	Data                []float64 // Y-values of the graph
}

// Draws the graph on the screen
func (g *Graph) Draw(screen *ebiten.Image) {
	// Draw horizontal grid lines
	for i := float32(0.0); i <= 1.0; i += 0.25 {
		vector.StrokeLine(screen, g.x, g.y+(g.height*i), g.x+g.width, g.y+(g.height*i), 1.0, colour.LightGray, false)
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
	step := g.width / float32(math.Max(float64(pointCount-1), 8))

	// Draw data lines
	for i := 0; i < pointCount-1; i++ {
		x1 := g.x + step*float32(i)
		x2 := g.x + step*float32(i+1)

		// Scale Y values to fit the graph
		y1 := g.y + g.height - float32((g.Data[i]-minValue)/valueRange)*g.height
		y2 := g.y + g.height - float32((g.Data[i+1]-minValue)/valueRange)*g.height

		vector.StrokeLine(screen, x1, y1, x2, y2, 2.0, colour.Green, true)
	}

	// add value label
	label := ""
	switch g.graphType {
	case Int:
		label = fmt.Sprintf("%.0f", lastValue)
	case Float:
		label = fmt.Sprintf("%.2f", lastValue)
	case Percentage:
		label = fmt.Sprintf("%.2f %%", lastValue)
	case Currency:
		label = utils.FormatCurrency(lastValue, "$")
	}
	ebitenutil.DebugPrintAt(screen, label, int(g.x)+8, int(g.y+g.height)-14)
}

func (g *Graph) Update() {}

func (g *Graph) SetOffset(x, y int) {
	g.x = float32(x)
	g.y = float32(y)
}
