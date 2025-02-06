package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

const (
	bottomBarHeight   = 24
	bottomButtonWidth = 36
)

type BottomBar struct {
	WindowsVisible            bool
	toggleWindows             func()
	screenHeight, screenWidth int
	bottomButtons             []*Button
}

func (b *BottomBar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, bottomButtonWidth, float32(b.screenHeight-bottomBarHeight), float32(b.screenWidth-bottomButtonWidth*2), bottomBarHeight, colour.DarkSemiBlack, true)
	ebitenutil.DebugPrintAt(screen, entities.Sim.GetStats(), bottomButtonWidth+10, b.screenHeight-bottomBarHeight+4)
	for i := range b.bottomButtons {
		b.bottomButtons[i].Draw(screen)
	}
}

func (b *BottomBar) Update() error {
	switch entities.Sim.SimulationSpeed {
	case entities.Slow:
		b.bottomButtons[0].Label = ">  "
	case entities.Mid:
		b.bottomButtons[0].Label = ">> "
	default:
		b.bottomButtons[0].Label = ">>>"
	}

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
		WindowsVisible: false,
		toggleWindows:  toggleWindows,
		screenHeight:   screenHeight,
		screenWidth:    screenWidth,
	}
	bar.bottomButtons = []*Button{
		{
			Label:      ">  ",
			X:          0,
			Y:          screenHeight - bottomBarHeight,
			Width:      bottomButtonWidth,
			Height:     bottomBarHeight,
			Color:      colour.DarkSemiBlack,
			HoverColor: colour.Blue,
			OnClick:    entities.Sim.ChangeSimulationSpeed,
		},
		{
			Label:      "[+]",
			X:          screenWidth - bottomButtonWidth,
			Y:          screenHeight - bottomBarHeight,
			Width:      bottomButtonWidth,
			Height:     bottomBarHeight,
			Color:      colour.DarkSemiBlack,
			HoverColor: colour.DarkGreen,
			OnClick:    toggleWindows,
		},
	}

	return bar
}
