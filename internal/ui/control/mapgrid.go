package control

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/entities"
)

const Border = 10

type MapGrid struct {
	x, y, size, gridSize, dataMax, frameCounter int
	data                                        [][]int
	dataSource                                  func() ([][]int, int) // Function to dynamically fetch data given an x, y coordinate
}

func (mg *MapGrid) getColour(x, y int) color.Color {
	if len(mg.data) > x && len(mg.data[x]) > y && mg.dataMax > 0 {
		intensity := uint8(255 * mg.data[x][y] / mg.dataMax)
		return color.RGBA{0, 63, 255, intensity}
	}
	return color.Transparent
}

func (mg *MapGrid) Draw(screen *ebiten.Image) {
	gridWidth := int(float32(mg.size-2*Border) * 0.707)
	grid := ebiten.NewImage(gridWidth, gridWidth)
	vector.DrawFilledRect(grid, 0, 0, float32(gridWidth), float32(gridWidth), color.White, false)
	cellSize := gridWidth / mg.gridSize
	for x := 0; x < mg.gridSize; x++ {
		for y := 0; y < mg.gridSize; y++ {
			vector.DrawFilledRect(grid, float32(x*cellSize), float32(y*cellSize),
				float32(cellSize), float32(cellSize), mg.getColour(x, y), false)
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(math.Pi / 4)
	op.GeoM.Scale(1, 0.5)
	op.GeoM.Translate(float64(mg.x+((mg.size+Border)/2)), float64(mg.y+Border))
	screen.DrawImage(grid, op)
}

func (mg *MapGrid) Update() {
	mg.frameCounter++
	if mg.frameCounter >= 60 { // update every second
		mg.frameCounter = 0
		entities.Sim.Mutex.RLock()
		mg.data, mg.dataMax = mg.dataSource()
		entities.Sim.Mutex.RUnlock()
	}
}

func (mg *MapGrid) SetOffset(x, y int) {
	mg.x = x
	mg.y = y
}

// NewMapGrid creates a new Map Grid
func NewMapGrid(x, y, size, gridSize int, dataSource func() ([][]int, int)) *MapGrid {
	return &MapGrid{
		x:          x,
		y:          y,
		size:       size,
		gridSize:   gridSize,
		dataSource: dataSource,
	}
}
