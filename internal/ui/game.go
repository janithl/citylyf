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
	"github.com/janithl/citylyf/utils"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

type Game struct {
	graph []Graph
}

func (g *Game) Update() error {
	g.graph[0].Data = entities.Sim.Market.MarketValues
	g.graph[0].MaxValue = entities.Sim.Market.MarketHigh

	g.graph[1].Data = utils.ConvertToF64(entities.Sim.People.PopulationValues)
	g.graph[1].MaxValue = float64(entities.Sim.People.PopulationHigh)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(image.White)
	for i := 0; i < len(g.graph); i++ {
		g.graph[i].Draw(screen)
	}

	vector.DrawFilledRect(screen, 0, screenHeight-24, screenWidth, 24, colour.Black, true)
	marketStats := fmt.Sprintf("%s | Population: %d (%.2f%%) | Unemployment: %.2f%% | Companies: %d | Market Value: %.2f (%.2f%%) | Inflation: %.2f%%",
		entities.Sim.Date.Format("2006-01-02"), entities.Sim.People.Population, entities.Sim.People.PopulationGrowthRate(),
		entities.Sim.People.UnemploymentRate(), len(entities.Sim.Companies), entities.Sim.Market.GetMarketValue(),
		entities.Sim.Market.LastMarketGrowthRate, entities.Sim.Market.LastInflationRate)
	ebitenutil.DebugPrintAt(screen, marketStats, 10, screenHeight-20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunGame() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("citylyf")
	game := &Game{
		graph: []Graph{
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
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
