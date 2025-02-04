package ui

import (
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
	screenWidth       = 1024
	screenHeight      = 768
	bottomBarHeight   = 24
	bottomButtonWidth = 36
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

	vector.DrawFilledRect(screen, bottomButtonWidth, screenHeight-bottomBarHeight, screenWidth, bottomBarHeight, colour.Black, true)
	ebitenutil.DebugPrintAt(screen, entities.Sim.GetStats(), bottomButtonWidth+10, screenHeight-bottomBarHeight+4)
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
				Width:      bottomButtonWidth,
				Height:     bottomBarHeight,
				Color:      colour.Black,
				HoverColor: colour.Blue,
				OnClick:    entities.Sim.ChangeSimulationSpeed,
			},
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
