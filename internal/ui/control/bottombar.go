package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type BottomBar struct {
	Enabled, WindowsVisible   bool
	toggleWindows             func()
	screenHeight, screenWidth int
	bottomButtons             []*Button
}

func (b *BottomBar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, buttonWidth, float32(b.screenHeight-buttonHeight), float32(b.screenWidth-buttonWidth*2), buttonHeight, colour.DarkSemiBlack, false)
	text := "<<< Close \"Map Control\" to begin the simulation >>>"
	if b.Enabled {
		text = entities.Sim.GetStats()
	}
	ebitenutil.DebugPrintAt(screen, text, buttonWidth+10, b.screenHeight-buttonHeight+4)

	for i := range b.bottomButtons {
		b.bottomButtons[i].Draw(screen)
	}
}

func (b *BottomBar) Update() error {
	buttonColour := colour.DarkSemiBlack
	switch entities.Sim.SimulationSpeed {
	case entities.Slow:
		b.bottomButtons[0].Label = ">  "
	case entities.Mid:
		b.bottomButtons[0].Label = ">> "
	case entities.Fast:
		b.bottomButtons[0].Label = ">>>"
	default:
		b.bottomButtons[0].Label = "|| "
		buttonColour = colour.DarkRed
	}
	b.bottomButtons[0].Color = buttonColour

	if b.WindowsVisible {
		b.bottomButtons[1].Label = "[-]"
	} else {
		b.bottomButtons[1].Label = "[+]"
	}

	for i := range b.bottomButtons {
		b.bottomButtons[i].Update()
	}

	return nil
}

func NewBottomBar(screenHeight, screenWidth int, toggleWindows func()) *BottomBar {
	bar := &BottomBar{
		Enabled:        false,
		WindowsVisible: false,
		toggleWindows:  toggleWindows,
		screenHeight:   screenHeight,
		screenWidth:    screenWidth,
	}
	bar.bottomButtons = []*Button{
		{
			Label:      ">  ",
			X:          0,
			Y:          screenHeight - buttonHeight,
			Width:      buttonWidth,
			Height:     buttonHeight,
			Color:      colour.DarkSemiBlack,
			HoverColor: colour.Blue,
			OnClick: func() {
				if bar.Enabled {
					entities.Sim.ChangeSimulationSpeed()
				}
			},
		},
		{
			Label:      "[+]",
			X:          screenWidth - buttonWidth,
			Y:          screenHeight - buttonHeight,
			Width:      buttonWidth,
			Height:     buttonHeight,
			Color:      colour.DarkSemiBlack,
			HoverColor: colour.DarkGreen,
			OnClick: func() {
				if bar.Enabled {
					toggleWindows()
				}
			},
		},
	}

	return bar
}
