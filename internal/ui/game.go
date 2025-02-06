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
	windows        []control.Window
	graphWindows   []control.GraphWindow
	windowsVisible bool
	bottomBar      control.BottomBar
	animations     []Animation
}

func (g *Game) Update() error {
	// replace the children of the stats window
	g.windows[0].ClearChildren()
	companies := []string{}
	for _, c := range entities.Sim.Companies {
		companies = append(companies, c.GetStats())
	}
	g.windows[0].AddChild(control.NewTextList(0, 0, 432, 360, companies))

	// replace the children of the households window
	g.windows[1].ClearChildren()
	households := []string{}
	for _, h := range entities.Sim.People.Households {
		households = append(households, h.GetStats())
	}
	g.windows[1].AddChild(control.NewTextList(0, 0, 360, 480, households))

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

	for i := range g.animations {
		g.animations[i].Update()
	}

	g.bottomBar.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colour.Gray)

	for i := range g.animations {
		g.animations[i].Draw(screen)
	}

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

func (g *Game) toggleAllWindows() {
	g.windowsVisible = !g.windowsVisible
	for i := range g.windows {
		g.windows[i].IsVisible = g.windowsVisible
	}
	for i := range g.graphWindows {
		g.graphWindows[i].Window.IsVisible = g.windowsVisible
	}
	g.bottomBar.WindowsVisible = g.windowsVisible
}

func RunGame() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("citylyf")

	game := &Game{
		windowsVisible: false,
		windows:        []control.Window{},
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
	}

	game.bottomBar = *control.NewBottomBar(screenHeight, screenWidth, game.toggleAllWindows)

	closeWindows := func(title string) {
		for i := range game.windows {
			if game.windows[i].Title == title {
				game.windows[i].CloseWindow()
			}
		}
	}

	companyWin := control.NewWindow(10, 290, 432, 360, "Companies", closeWindows)
	companies := []string{}
	for _, c := range entities.Sim.Companies {
		companies = append(companies, c.GetStats())
	}
	companyWin.AddChild(control.NewTextList(0, 0, 432, 360, companies))

	householdsWin := control.NewWindow(640, 10, 360, 480, "Households", closeWindows)
	households := []string{}
	for _, h := range entities.Sim.People.Households {
		households = append(households, h.GetStats())
	}
	householdsWin.AddChild(control.NewTextList(0, 0, 360, 480, households))

	popGraphWin := control.NewWindow(10, 10, 200, 130, "Population", closeWindows)
	popGraphWin.AddChild(&control.Graph{
		X:      0,
		Y:      2,
		Width:  200,
		Height: 100,
		Data:   utils.ConvertToF64(entities.Sim.People.PopulationValues),
	})

	game.windows = append(game.windows, *companyWin, *householdsWin, *popGraphWin)

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
			Window: control.NewWindow(10, 150, 160, 100, "Market Growth Rate", game.closeGraphWindows),
		},
		{
			Data:   &entities.Sim.Market.History.MarketSentiment,
			Window: control.NewWindow(180, 150, 160, 100, "Market Sentiment", game.closeGraphWindows),
		},
		{
			Data:   &entities.Sim.Market.History.CompanyProfits,
			Window: control.NewWindow(350, 150, 160, 100, "Company Profits", game.closeGraphWindows),
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
