package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/ui/colour"
	"github.com/janithl/citylyf/internal/utils"
)

type MainMenu struct {
	x, y, width, height, screenWidth, screenHeight, startIndex int
	entries                                                    []string
	layoutGrid                                                 *Grid
	onEntryClick                                               func(*string)
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, float32(m.screenWidth), float32(m.screenHeight), colour.DarkSemiBlack, false)
	m.layoutGrid.Draw(screen)
}

func (m *MainMenu) Update() {
	// Mouse wheel zoom
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) && m.startIndex < len(m.entries)-(m.layoutGrid.rows-2) {
		m.startIndex += 1 // Scroll down
		m.updateEntries()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) && m.startIndex > 0 {
		m.startIndex -= 1 // Scroll up
		m.updateEntries()
	}

	m.layoutGrid.Update()
}

func (m *MainMenu) updateEntries() {
	savedir := utils.GetSaveDir()
	for i, entry := range m.entries[m.startIndex:] {
		if i >= m.layoutGrid.rows-2 {
			break
		}

		path := savedir + "/" + entry
		if m.layoutGrid.Children[i+2][0] != nil {
			m.layoutGrid.Children[i+2][0].(*Button).Label = entry
			m.layoutGrid.Children[i+2][0].(*Button).OnClick = func() { m.onEntryClick(&path) }
		}
	}
}

func (m *MainMenu) Layout(width, height int) {
	m.screenWidth = width
	m.screenHeight = height
	m.x = (width - m.width*3) / 2
	m.y = (height - m.height) / 2
	m.layoutGrid.SetOffset(m.x, m.y)
}

func NewMainMenu(width, maxEntries int, resumable bool, toggleMenuMode, loadGame, endGame func(), startNewGame func(*string)) *MainMenu {
	menu := &MainMenu{
		x:          0,
		y:          0,
		width:      width,
		height:     maxEntries * menuEntryHeight,
		layoutGrid: NewGrid(0, 0, width, maxEntries*menuEntryHeight, 1, maxEntries),
	}

	if resumable {
		menu.layoutGrid.Children[0][0] = &Button{Label: "Resume Game", Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: toggleMenuMode}
	}
	menu.layoutGrid.Children[1][0] = &Button{Label: "New Game", Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: func() { startNewGame(nil) }}
	menu.layoutGrid.Children[2][0] = &Button{Label: "Load Game", Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: loadGame}
	menu.layoutGrid.Children[3][0] = &Button{Label: "Exit", Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: endGame}

	return menu
}

func NewLoadGameMenu(width, maxEntries int, loadMainMenu func(), startNewGame func(*string)) *MainMenu {
	menu := &MainMenu{
		x:            0,
		y:            0,
		width:        width,
		height:       maxEntries * menuEntryHeight,
		layoutGrid:   NewGrid(0, 0, width, maxEntries*menuEntryHeight, 1, maxEntries),
		onEntryClick: startNewGame,
	}

	menu.layoutGrid.Children[0][0] = &Button{Label: "<- Back", Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: loadMainMenu}

	for i := 2; i < maxEntries; i++ {
		menu.layoutGrid.Children[i][0] = &Button{Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red}
	}

	if savedir := utils.GetSaveDir(); savedir != "" {
		if files, err := utils.GetDirFiles(savedir); err == nil {
			menu.entries = files
			menu.updateEntries()
		}
	}

	return menu
}
