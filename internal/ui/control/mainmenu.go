package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/gamefile"
	"github.com/janithl/citylyf/internal/ui/colour"
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
	_, scrollY := ebiten.Wheel() // Mouse wheel scroll
	scrollUp := inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) || scrollY < 0
	scrollDown := inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) || scrollY > 0

	if scrollDown && m.startIndex < len(m.entries)-(m.layoutGrid.rows-2) {
		m.startIndex += 1
		m.updateEntries()
	} else if scrollUp && m.startIndex > 0 {
		m.startIndex -= 1
		m.updateEntries()
	}

	m.layoutGrid.Update()
}

func (m *MainMenu) updateEntries() {
	savedir := gamefile.GetSavesDir()
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
	exitBtn := &Button{Label: "Exit", Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: endGame}
	if entities.Sim != nil && entities.Sim.SavePath != "" {
		menu.layoutGrid.Children[3][0] = &Button{Label: "Save Game", Width: width, Height: buttonHeight, Scale: 3, Color: colour.Transparent, HoverColor: colour.Red, OnClick: func() { gamefile.Save(entities.Sim.SavePath); toggleMenuMode() }}
		menu.layoutGrid.Children[4][0] = exitBtn
	} else {
		menu.layoutGrid.Children[3][0] = exitBtn
	}

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

	if savesdir := gamefile.GetSavesDir(); savesdir != "" {
		menu.entries = gamefile.GetDirFiles(savesdir)
		menu.updateEntries()
	}

	return menu
}
