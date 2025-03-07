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
	bottomText                string
}

func (b *BottomBar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, buttonWidth, float32(b.screenHeight-buttonHeight), float32(b.screenWidth-buttonWidth*2), buttonHeight, colour.DarkSemiBlack, false)
	ebitenutil.DebugPrintAt(screen, b.bottomText, buttonWidth+10, b.screenHeight-buttonHeight+4)

	for i := range b.bottomButtons {
		b.bottomButtons[i].Draw(screen)
	}
}

func (b *BottomBar) Update() error {
	if b.Enabled {
		select { // non-blocking read from stats channel
		case stats := <-entities.SimStats:
			b.bottomText = stats
		default:
		}
	}

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
		bottomText:     "<<< Close \"Map Control\" to begin the simulation >>>",
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
					entities.Sim.Mutex.Lock()
					entities.Sim.ChangeSimulationSpeed()
					entities.Sim.Mutex.Unlock()
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
