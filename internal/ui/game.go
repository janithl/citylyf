package ui

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

type Game struct {
	graph Graph
}

func (g *Game) Update() error {
	g.graph.Data = entities.Sim.Market.MarketValues
	g.graph.MaxValue = entities.Sim.Market.MarketHigh
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(image.White)
	g.graph.Draw(screen)

	vector.DrawFilledRect(screen, 0, screenHeight-24, screenWidth, 24, colour.Black, true)
	marketStats := fmt.Sprintf("%s | Population: %d | Unemployment: %.2f%% | Market Value: %.2f (%.2f%%) | Inflation: %.2f%%",
		entities.Sim.Date.Format("2006-01-02"), entities.Sim.People.Population, entities.Sim.People.UnemploymentRate(),
		entities.Sim.Market.GetMarketValue(), entities.Sim.Market.LastMarketGrowthRate, entities.Sim.Market.LastInflationRate)
	ebitenutil.DebugPrintAt(screen, marketStats, 10, screenHeight-20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunGame() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("citylyf")
	game := &Game{
		graph: Graph{
			X:        0,
			Y:        0,
			Width:    200,
			Height:   100,
			Data:     entities.Sim.Market.MarketValues,
			MaxValue: entities.Sim.Market.MarketHigh,
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
