package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
	"github.com/janithl/citylyf/internal/ui/control"
	"github.com/janithl/citylyf/internal/utils"
)

type WindowSystem struct {
	windowsVisible, generatingMap bool
	windows                       []control.Window
	listWindows                   []control.ListWindow
	graphWindows                  []control.GraphWindow
	bottomBar                     *control.BottomBar
	generateBar                   []control.Button
}

func (ws *WindowSystem) Update() error {
	for i := range ws.windows {
		ws.windows[i].Update()
	}
	for i := range ws.listWindows {
		ws.listWindows[i].Update()
	}
	for i := range ws.graphWindows {
		ws.graphWindows[i].Update()
	}
	ws.bottomBar.Update()
	for i := range ws.generateBar {
		ws.generateBar[i].Update()
	}
	return nil
}

func (ws *WindowSystem) Draw(screen *ebiten.Image) {
	if ws.generatingMap {
		for i := range ws.generateBar {
			ws.generateBar[i].Draw(screen)
		}
	}

	for i := range ws.windows {
		ws.windows[i].Draw(screen)
	}
	for i := range ws.listWindows {
		ws.listWindows[i].Draw(screen)
	}
	for i := range ws.graphWindows {
		ws.graphWindows[i].Draw(screen)
	}
	ws.bottomBar.Draw(screen)
}

func (ws *WindowSystem) closeWindows(title string) {
	for i := range ws.windows {
		if ws.windows[i].Title == title {
			ws.windows[i].CloseWindow()
			return
		}
	}
	for i := range ws.listWindows {
		if ws.listWindows[i].Window.Title == title {
			ws.listWindows[i].Window.CloseWindow()
			return
		}
	}
	for i := range ws.graphWindows {
		if ws.graphWindows[i].Window.Title == title {
			ws.graphWindows[i].Window.CloseWindow()
			return
		}
	}
}

func (ws *WindowSystem) toggleAllWindows() {
	ws.windowsVisible = !ws.windowsVisible
	for i := range ws.windows {
		ws.windows[i].IsVisible = ws.windowsVisible
	}
	for i := range ws.listWindows {
		ws.listWindows[i].Window.IsVisible = ws.windowsVisible
	}
	for i := range ws.graphWindows {
		ws.graphWindows[i].Window.IsVisible = ws.windowsVisible
	}
	ws.bottomBar.WindowsVisible = ws.windowsVisible
}

func (ws *WindowSystem) onWindowItemClick(title string, index int) {
	fmt.Println("Learn more about", title, "#", index)
}

func (ws *WindowSystem) doneGeneratingMap() {
	ws.generatingMap = false
}

func NewWindowSystem() *WindowSystem {
	ws := &WindowSystem{
		windowsVisible: false,
		generatingMap:  true,
		windows:        []control.Window{},
	}

	ppWin := *control.NewWindow(850, 10, 360, 270, "Population Pyramid", ws.closeWindows)
	ppWin.AddChild(&control.PopulationPyramid{X: 0, Y: 0, Width: 360, Height: 250})
	ws.windows = append(ws.windows, ppWin)

	ws.listWindows = []control.ListWindow{
		*control.NewListWindow(10, 290, 500, 360, "Companies", ws.closeWindows, ws.onWindowItemClick,
			func() []string {
				companies := []string{}
				for _, id := range entities.Sim.CompanyIDs {
					companies = append(companies, entities.Sim.Companies[id].GetStats())
				}

				return companies
			}),
		*control.NewListWindow(520, 290, 460, 360, "Households", ws.closeWindows, ws.onWindowItemClick,
			func() []string {
				households := []string{}
				for _, h := range entities.Sim.People.Households {
					households = append(households, h.GetStats())
				}
				return households
			}),
	}

	ws.graphWindows = []control.GraphWindow{
		*control.NewGraphWindow(10, 10, 200, 130, "Population", ws.closeWindows, control.Int,
			func() []float64 { return utils.ConvertToF64(entities.Sim.People.PopulationValues) }),
		*control.NewGraphWindow(220, 10, 200, 130, "Market Value", ws.closeWindows, control.Float,
			func() []float64 { return entities.Sim.Market.History.MarketValue }),
		*control.NewGraphWindow(430, 10, 200, 130, "Inflation Rate", ws.closeWindows, control.Percentage,
			func() []float64 { return entities.Sim.Market.History.InflationRate }),
		*control.NewGraphWindow(640, 10, 200, 130, "Gov Reserves", ws.closeWindows, control.Currency,
			func() []float64 { return utils.ConvertToF64(entities.Sim.Government.ReserveValues) }),
		*control.NewGraphWindow(10, 150, 160, 100, "Market Growth Rate", ws.closeWindows, control.Percentage,
			func() []float64 { return entities.Sim.Market.History.MarketGrowthRate }),
		*control.NewGraphWindow(180, 150, 160, 100, "Market Sentiment", ws.closeWindows, control.Float,
			func() []float64 { return entities.Sim.Market.History.MarketSentiment }),
		*control.NewGraphWindow(350, 150, 160, 100, "Company Profits", ws.closeWindows, control.Currency,
			func() []float64 { return entities.Sim.Market.History.CompanyProfits }),
		*control.NewGraphWindow(520, 150, 155, 100, "Collected Tax", ws.closeWindows, control.Currency,
			func() []float64 { return utils.ConvertToF64(entities.Sim.Government.CollectedTaxValues) }),
		*control.NewGraphWindow(685, 150, 155, 100, "Unemployment Rate", ws.closeWindows, control.Percentage,
			func() []float64 { return entities.Sim.People.UnemploymentRateValues }),
	}

	ws.bottomBar = control.NewBottomBar(screenHeight, screenWidth, ws.toggleAllWindows)
	ws.generateBar = []control.Button{
		{Label: "Regenerate Map", X: 4, Y: 4, Width: 200, Height: 24, Color: colour.Black, HoverColor: colour.DarkGreen, OnClick: entities.Sim.RegenerateMap},
		{Label: "Done", X: 4, Y: 32, Width: 200, Height: 24, Color: colour.Black, HoverColor: colour.DarkGreen, OnClick: ws.doneGeneratingMap},
	}

	return ws
}
