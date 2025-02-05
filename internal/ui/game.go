package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
	"github.com/janithl/citylyf/internal/ui/control"
	"github.com/janithl/citylyf/internal/utils"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

type Game struct {
	windows      []control.Window
	graphWindows []control.GraphWindow
	bottomBar    control.BottomBar
}

func (g *Game) Update() error {
	// replace the children of the stats window
	g.windows[0].ClearChildren()
	for i := range entities.Sim.Companies {
		label := &control.Label{X: 6, Y: 4 + (i * 16), Text: entities.Sim.Companies[i].GetStats()}
		g.windows[0].AddChild(label)
	}

	// replace the children of the households window
	g.windows[1].ClearChildren()
	for j := range entities.Sim.People.Households {
		label := &control.Label{X: 6, Y: 4 + (j * 16), Text: entities.Sim.People.Households[j].GetStats()}
		g.windows[1].AddChild(label)
	}

	// replace the children of the population graph window
	g.windows[2].ClearChildren()
	g.windows[2].AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   utils.ConvertToF64(entities.Sim.People.PopulationValues),
	})

	for i := range g.windows {
		g.windows[i].Update()
	}

	for i := range g.graphWindows {
		g.graphWindows[i].Update()
	}

	g.bottomBar.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colour.Gray)

	for i := range g.windows {
		g.windows[i].Draw(screen)
	}

	for i := range g.graphWindows {
		g.graphWindows[i].Draw(screen)
	}

	g.bottomBar.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) closeGraphWindows(title string) {
	for i := range g.graphWindows {
		if g.graphWindows[i].Window.Title == title {
			g.graphWindows[i].Window.CloseWindow()
		}
	}
}

func RunGame() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("citylyf")

	game := &Game{
		bottomBar: *control.NewBottomBar(screenHeight, screenWidth),
		windows:   []control.Window{},
	}

	closeWindows := func(title string) {
		for i := range game.windows {
			if game.windows[i].Title == title {
				game.windows[i].CloseWindow()
			}
		}
	}

	statsWin := control.NewWindow(10, 290, 432, 360, "Company Stats", closeWindows)
	for i := range entities.Sim.Companies {
		label := &control.Label{X: 6, Y: 4 + (i * 16), Text: entities.Sim.Companies[i].GetStats()}
		statsWin.AddChild(label)
	}

	householdsWin := control.NewWindow(640, 10, 360, 600, "Households", closeWindows)
	for j := range entities.Sim.People.Households {
		label := &control.Label{X: 6, Y: 4 + (j * 16), Text: entities.Sim.People.Households[j].GetStats()}
		householdsWin.AddChild(label)
	}

	popGraphWin := control.NewWindow(10, 10, 200, 130, "Population", closeWindows)
	popGraphWin.AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   utils.ConvertToF64(entities.Sim.People.PopulationValues),
	})

	game.windows = append(game.windows, *statsWin, *householdsWin, *popGraphWin)

	game.graphWindows = []control.GraphWindow{
		{
			Data:   &entities.Sim.Market.History.MarketValue,
			Window: control.NewWindow(220, 10, 200, 130, "Market Value", game.closeGraphWindows),
		},
		{
			Data:   &entities.Sim.Market.History.InflationRate,
			Window: control.NewWindow(430, 10, 200, 130, "Inflation Rate", game.closeGraphWindows),
		},
		{
			Data:   &entities.Sim.Market.History.MarketGrowthRate,
			Window: control.NewWindow(10, 150, 200, 130, "Market Growth Rate", game.closeGraphWindows),
		},
		{
			Data:   &entities.Sim.Market.History.MarketSentiment,
			Window: control.NewWindow(220, 150, 200, 130, "Market Sentiment", game.closeGraphWindows),
		},
		{
			Data:   &entities.Sim.Market.History.CompanyProfits,
			Window: control.NewWindow(430, 150, 200, 130, "Company Profits", game.closeGraphWindows),
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
