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
	windows   []control.Window
	bottomBar control.BottomBar
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

	// replace the children of the graphs windows
	g.windows[2].ClearChildren()
	g.windows[3].ClearChildren()
	g.windows[4].ClearChildren()
	g.windows[5].ClearChildren()
	g.windows[6].ClearChildren()
	g.windows[7].ClearChildren()
	g.windows[2].AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   utils.ConvertToF64(entities.Sim.People.PopulationValues),
	})
	g.windows[3].AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.MarketValue,
	})
	g.windows[4].AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.InflationRate,
	})
	g.windows[5].AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.MarketGrowthRate,
	})
	g.windows[6].AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.MarketSentiment,
	})
	g.windows[7].AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.CompanyProfits,
	})

	for i := range g.windows {
		g.windows[i].Update()
	}

	g.bottomBar.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colour.Gray)

	for k := range g.windows {
		g.windows[k].Draw(screen)
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

	marketGraphWin := control.NewWindow(220, 10, 200, 130, "Market Value", closeWindows)
	marketGraphWin.AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.MarketValue,
	})

	inflationGraphWin := control.NewWindow(430, 10, 200, 130, "Inflation", closeWindows)
	inflationGraphWin.AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.InflationRate,
	})

	growthGraphWin := control.NewWindow(10, 150, 200, 130, "Market Growth", closeWindows)
	growthGraphWin.AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.MarketGrowthRate,
	})

	sentimentGraphWin := control.NewWindow(220, 150, 200, 130, "Market Sentiment", closeWindows)
	sentimentGraphWin.AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.MarketSentiment,
	})

	profitsGraphWin := control.NewWindow(430, 150, 200, 130, "Company Profits", closeWindows)
	profitsGraphWin.AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   entities.Sim.Market.History.CompanyProfits,
	})

	game.windows = append(game.windows, *statsWin, *householdsWin, *popGraphWin, *marketGraphWin,
		*inflationGraphWin, *growthGraphWin, *sentimentGraphWin, *profitsGraphWin)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
