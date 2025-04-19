package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type MainMenu struct {
	x, y, width, height, screenWidth, screenHeight int
	layoutGrid                                     *Grid
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, float32(m.screenWidth), float32(m.screenHeight), colour.DarkSemiBlack, false)
	m.layoutGrid.Draw(screen)
}

func (m *MainMenu) Update() {
	m.layoutGrid.Update()
}

func (m *MainMenu) Layout(width, height int) {
	m.screenWidth = width
	m.screenHeight = height
	m.x = (width - m.width*3) / 2
	m.y = (height - m.height) / 2
	m.layoutGrid.SetOffset(m.x, m.y)
}

func NewMainMenu(width, height int, toggleMenuMode, endGame func()) *MainMenu {
	menu := &MainMenu{
		x:          0,
		y:          0,
		width:      width,
		height:     height,
		layoutGrid: NewGrid(0, 0, width, height, 1, 4),
	}

	menu.layoutGrid.Children[0][0] = &Button{Label: "Resume Game", X: 0, Y: 0, Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: toggleMenuMode}
	menu.layoutGrid.Children[1][0] = &Button{Label: "New Game", X: 0, Y: 0, Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: func() {}}
	menu.layoutGrid.Children[2][0] = &Button{Label: "Load Game", X: 0, Y: 0, Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: func() {}}
	menu.layoutGrid.Children[3][0] = &Button{Label: "Exit", X: 0, Y: 0, Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: endGame}

	return menu
}
