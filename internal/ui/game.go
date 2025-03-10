package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/ui/colour"
	"github.com/janithl/citylyf/internal/ui/world"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

type Game struct {
	worldRenderer world.WorldRenderer
	windowSystem  WindowSystem
}

func (g *Game) Update() error {
	g.worldRenderer.Update()
	g.windowSystem.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colour.Gray)
	g.worldRenderer.Draw(screen)
	g.windowSystem.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunGame() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("citylyf")

	game := &Game{
		worldRenderer: *world.NewWorldRenderer(screenWidth, screenHeight),
		windowSystem:  *NewWindowSystem(),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
