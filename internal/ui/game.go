package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
	"github.com/janithl/citylyf/internal/ui/control"
	"github.com/janithl/citylyf/utils"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

type Game struct {
	graphs    []control.Graph
	buttons   []control.Button
	bottomBar control.BottomBar
}

func (g *Game) Update() error {
	g.graphs[0].Data = entities.Sim.Market.MarketValues
	g.graphs[0].MaxValue = entities.Sim.Market.MarketHigh

	g.graphs[1].Data = utils.ConvertToF64(entities.Sim.People.PopulationValues)
	g.graphs[1].MaxValue = float64(entities.Sim.People.PopulationHigh)

	g.bottomBar.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colour.Gray)
	for i := range g.graphs {
		g.graphs[i].Draw(screen)
	}

	for j := range g.buttons {
		g.buttons[j].Draw(screen)
	}

	g.bottomBar.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunGame() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("citylyf")
	game := &Game{
		graphs: []control.Graph{
			{
				Title:    "Market Value",
				X:        0,
				Y:        0,
				Width:    200,
				Height:   100,
				Data:     entities.Sim.Market.MarketValues,
				MaxValue: entities.Sim.Market.MarketHigh,
			}, {
				Title:    "Population",
				X:        210,
				Y:        0,
				Width:    200,
				Height:   100,
				Data:     utils.ConvertToF64(entities.Sim.People.PopulationValues),
				MaxValue: float64(entities.Sim.People.PopulationHigh),
			},
		},
		bottomBar: *control.NewBottomBar(screenHeight, screenWidth),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
