package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/ui/colour"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

type Game struct {
	animations   []Animation
	windowSystem WindowSystem
}

func (g *Game) Update() error {
	for i := range g.animations {
		g.animations[i].Update()
	}

	g.windowSystem.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colour.Gray)

	for i := range g.animations {
		g.animations[i].Draw(screen)
	}

	g.windowSystem.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunGame() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("citylyf")

	game := &Game{
		animations: []Animation{
			*NewAnimation(screenWidth/2, screenHeight/2, 0.4, 0),
			*NewAnimation(screenWidth/2, screenHeight/2, 0.28, 0.28),
			*NewAnimation(screenWidth/2, screenHeight/2, 0.28, -0.28),

			*NewAnimation(screenWidth/2, screenHeight/2, 0, 0.4),
			*NewAnimation(screenWidth/2, screenHeight/2, 0, -0.4),
			*NewAnimation(screenWidth/2, screenHeight/2, 0, 0),

			*NewAnimation(screenWidth/2, screenHeight/2, -0.28, -0.28),
			*NewAnimation(screenWidth/2, screenHeight/2, -0.28, 0.28),
			*NewAnimation(screenWidth/2, screenHeight/2, -0.4, 0),
		},
		windowSystem: *NewWindowSystem(),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
