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
)

type Game struct {
	worldRenderer *world.WorldRenderer
	windowSystem  *WindowSystem
	mapRegenMode  bool
	mapControl    *control.MapControl
}

func (g *Game) Update() error {
	g.worldRenderer.Update(g.mapRegenMode)
	g.windowSystem.Update()
	g.mapControl.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colour.Gray)
	g.worldRenderer.Draw(screen)
	if g.mapRegenMode {
		g.mapControl.Draw(screen)
	} else {
		g.windowSystem.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) EndRegenMode() {
	g.mapRegenMode = false
	entities.Sim.Mutex.Lock()
	entities.Sim.ChangeSimulationSpeed()
	entities.Sim.Mutex.Unlock()
}

func RunGame() {
	ebiten.SetWindowSize(screenWidth, screenHeight)

	gameTitle := "citylyf"
	if entities.Sim.SavePath != "" {
		gameTitle = entities.Sim.SavePath + " — " + gameTitle
	}
	ebiten.SetWindowTitle(gameTitle)

	game := &Game{
		worldRenderer: world.NewWorldRenderer(screenWidth, screenHeight),
		windowSystem:  NewWindowSystem(),
		mapRegenMode:  entities.Sim.SavePath == "",
	}
	game.mapControl = control.NewMapControl(0, 0, 240, 140, game.EndRegenMode)
	game.mapControl.SetOffset(screenWidth-240, screenHeight-140)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
