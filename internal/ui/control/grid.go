package control

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Grid struct {
	x, y, width, height, columns, rows int
	Children                           [][]Renderable // Child controls
}

func (g *Grid) Draw(screen *ebiten.Image) {
	// draw child controls
	for _, rowCells := range g.Children {
		for _, cell := range rowCells {
			if cell != nil {
				cell.Draw(screen)
			}
		}
	}
}

func (g *Grid) Update() {
	rowHeight := g.height / g.rows
	colWidth := g.width / g.columns

	// position and update child controls
	for rowNum, rowCells := range g.Children {
		for colNum, cell := range rowCells {
			if cell != nil {
				cell.SetOffset(g.x+colWidth*colNum, g.y+rowHeight*rowNum)
				cell.Update()
			}
		}
	}
}

func (g *Grid) SetOffset(x, y int) {
	g.x = x
	g.y = y
}

func NewGrid(x, y, width, height, columns, rows int) *Grid {
	children := make([][]Renderable, rows)
	for i := range rows {
		children[i] = make([]Renderable, columns)
	}

	return &Grid{
		x:        x,
		y:        y,
		width:    width,
		height:   height,
		rows:     rows,
		columns:  columns,
		Children: children,
	}
}
