package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/control"
	"github.com/janithl/citylyf/internal/utils"
)

type WindowSystem struct {
	windowsVisible bool
	windows        []control.Window
	listWindows    []control.ListWindow
	graphWindows   []control.GraphWindow
	bottomBar      *control.BottomBar
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
	return nil
}

func (ws *WindowSystem) Draw(screen *ebiten.Image) {
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
	switch title {
	case "Companies":
		company, exists := entities.Sim.Companies[index]
		if exists {
			fmt.Println(company.Name, company.CompanyAge(), company.Industry, company.GetNumberOfEmployees(), company.GetNumberOfJobOpenings())
			for _, emp := range company.GetEmployees() {
				fmt.Println(emp)
			}
		}
	case "Households":
		household, exists := entities.Sim.People.Households[index]
		if exists {
			fmt.Println(household.FamilyName(), household.HouseID, household.Size(), household.MoveInDate.Year())
			fmt.Println(household.GetMemberStats())
		}
	}
}

func NewWindowSystem() *WindowSystem {
	ws := &WindowSystem{
		windowsVisible: false,
		windows:        []control.Window{},
	}

	ppWin := *control.NewWindow(970, 10, 300, 270, "Population Pyramid", ws.closeWindows)
	ppWin.AddChild(&control.PopulationPyramid{X: 0, Y: 0, Width: 300, Height: 250})
	ws.windows = append(ws.windows, ppWin)

	gridWin := *control.NewWindow(990, 290, 240, 160, "Population Map", ws.closeWindows)
	gridWin.AddChild(control.NewMapGrid(0, 0, 240, 8, entities.Sim.Geography.Regions.GetPopulationStats))
	ws.windows = append(ws.windows, gridWin)

	ws.listWindows = []control.ListWindow{
		*control.NewListWindow(10, 290, 500, 360, "Companies", ws.closeWindows, ws.onWindowItemClick,
			func() []control.Statable {
				companies := []control.Statable{}
				for _, companyID := range entities.Sim.Companies.GetIDs() {
					company, exists := entities.Sim.Companies[companyID]
					if exists {
						companies = append(companies, company)
					}
				}
				return companies
			}),
		*control.NewListWindow(520, 290, 460, 360, "Households", ws.closeWindows, ws.onWindowItemClick,
			func() []control.Statable {
				households := []control.Statable{}
				for _, householdID := range entities.Sim.People.GetHouseholdIDs() {
					household, exists := entities.Sim.People.Households[householdID]
					if exists {
						households = append(households, household)
					}
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
		*control.NewGraphWindow(10, 150, 150, 100, "Market Growth Rate", ws.closeWindows, control.Percentage,
			func() []float64 { return entities.Sim.Market.History.MarketGrowthRate }),
		*control.NewGraphWindow(170, 150, 150, 100, "Market Sentiment", ws.closeWindows, control.Float,
			func() []float64 { return entities.Sim.Market.History.MarketSentiment }),
		*control.NewGraphWindow(330, 150, 150, 100, "Company Profits", ws.closeWindows, control.Currency,
			func() []float64 { return entities.Sim.Market.History.CompanyProfits }),
		*control.NewGraphWindow(490, 150, 150, 100, "Gov. Income", ws.closeWindows, control.Currency,
			func() []float64 { return utils.ConvertToF64(entities.Sim.Government.IncomeValues) }),
		*control.NewGraphWindow(650, 150, 150, 100, "Unemployment Rate", ws.closeWindows, control.Percentage,
			func() []float64 { return entities.Sim.People.UnemploymentRateValues }),
		*control.NewGraphWindow(810, 150, 150, 100, "Interest Rate", ws.closeWindows, control.Percentage,
			func() []float64 { return entities.Sim.Market.History.InterestRate }),
	}

	ws.bottomBar = control.NewBottomBar(screenHeight, screenWidth, ws.toggleAllWindows)
	return ws
}
