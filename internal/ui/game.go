package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
	"github.com/janithl/citylyf/internal/ui/control"
	"github.com/janithl/citylyf/internal/ui/world"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	mcWidth      = 256
	mcHeight     = 160
)

type Game struct {
	worldRenderer *world.WorldRenderer
	windowSystem  *WindowSystem
	mainMenu      *control.MainMenu
	mapControl    *control.MapControl
	startGame     func()

	terminate bool
}

func (g *Game) Update() error {
	if g.terminate {
		return ebiten.Termination
	}

	if g.mainMenu != nil {
		g.mainMenu.Update()
		return nil
	}

	if g.mapControl != nil {
		g.mapControl.Update()
	} else {
		g.worldRenderer.Update()
		g.windowSystem.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.mainMenu != nil {
		screen.Fill(colour.Black)
		g.mainMenu.Draw(screen)
		return
	}

	screen.Fill(colour.Gray)
	if g.worldRenderer != nil {
		g.worldRenderer.Draw(screen)
	}
	if g.mapControl != nil {
		g.mapControl.Draw(screen)
	} else {
		g.windowSystem.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.worldRenderer != nil {
		g.worldRenderer.Layout(outsideWidth, outsideHeight)
	}
	if g.windowSystem != nil {
		g.windowSystem.Layout(outsideWidth, outsideHeight)
	}
	if g.mainMenu != nil {
		g.mainMenu.Layout(outsideWidth, outsideHeight)
	}
	if g.mapControl != nil {
		g.mapControl.SetOffset(outsideWidth-mcWidth, outsideHeight-mcHeight)
	}

	return outsideWidth, outsideHeight
}

func (g *Game) EndRegenMode() {
	g.mapControl = nil

	entities.Sim.Mutex.Lock()
	entities.Sim.ChangeSimulationSpeed()
	entities.Sim.Mutex.Unlock()

	g.windowSystem = NewWindowSystem()
}

func (g *Game) EndGame() {
	g.terminate = true
}

func (g *Game) ToggleMenuMode() {
	if g.mainMenu != nil {
		g.mainMenu = nil
	} else {
		g.mainMenu = control.NewMainMenu(192, 288, true, g.ToggleMenuMode, g.StartNewGame, g.EndGame)
	}
}

func (g *Game) StartNewGame() {
	g.startGame()
	g.mainMenu = nil
	g.mapControl = control.NewMapControl(0, 0, mcWidth, mcHeight, g.EndRegenMode)
	g.mapControl.SetOffset(screenWidth-mcWidth, screenHeight-mcHeight)
	g.worldRenderer = world.NewWorldRenderer(screenWidth, screenHeight, g.ToggleMenuMode)
}

func RunGame(startGame func()) {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("citylyf")

	game := &Game{startGame: startGame}
	game.mainMenu = control.NewMainMenu(192, 288, false, game.ToggleMenuMode, game.StartNewGame, game.EndGame)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
