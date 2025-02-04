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
	"github.com/janithl/citylyf/internal/ui/control"
	"github.com/janithl/citylyf/utils"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

type Game struct {
	graphs  []control.Graph
	buttons []control.Button
}

func (g *Game) Update() error {
	g.graphs[0].Data = entities.Sim.Market.MarketValues
	g.graphs[0].MaxValue = entities.Sim.Market.MarketHigh

	g.graphs[1].Data = utils.ConvertToF64(entities.Sim.People.PopulationValues)
	g.graphs[1].MaxValue = float64(entities.Sim.People.PopulationHigh)

	switch entities.Sim.SimulationSpeed {
	case entities.Slow:
		g.buttons[0].Label = ">  "
	case entities.Mid:
		g.buttons[0].Label = ">> "
	default:
		g.buttons[0].Label = ">>>"
	}

	for i := range g.buttons {
		g.buttons[i].Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(image.White)
	for i := range g.graphs {
		g.graphs[i].Draw(screen)
	}

	for j := range g.buttons {
		g.buttons[j].Draw(screen)
	}

	vector.DrawFilledRect(screen, 36, screenHeight-24, screenWidth, 24, colour.Black, true)
	marketStats := fmt.Sprintf("%s | Population: %d (%.2f%%) | Unemployment: %.2f%% | Companies: %d | Market Value: %.2f (%.2f%%) | Inflation: %.2f%%",
		entities.Sim.Date.Format("2006-01-02"), entities.Sim.People.Population, entities.Sim.People.PopulationGrowthRate(),
		entities.Sim.People.UnemploymentRate(), len(entities.Sim.Companies), entities.Sim.Market.GetMarketValue(),
		entities.Sim.Market.LastMarketGrowthRate, entities.Sim.Market.LastInflationRate)
	ebitenutil.DebugPrintAt(screen, marketStats, 46, screenHeight-20)
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
		buttons: []control.Button{
			{
				Label:      ">  ",
				X:          0,
				Y:          screenHeight - 24,
				Width:      36,
				Height:     24,
				Color:      colour.Black,
				HoverColor: colour.Blue,
				OnClick: func() {
					switch entities.Sim.SimulationSpeed {
					case entities.Slow:
						entities.Sim.SimulationSpeed = entities.Mid
					case entities.Mid:
						entities.Sim.SimulationSpeed = entities.Fast
					default:
						entities.Sim.SimulationSpeed = entities.Slow
					}
				},
			},
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
