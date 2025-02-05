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
	screenHeight, screenWidth int
	bottomButton              Button
}

func (b *BottomBar) Draw(screen *ebiten.Image) {
	b.bottomButton.Draw(screen)
	vector.DrawFilledRect(screen, bottomButtonWidth, float32(b.screenHeight-bottomBarHeight), float32(b.screenWidth), bottomBarHeight, colour.Black, true)
	ebitenutil.DebugPrintAt(screen, entities.Sim.GetStats(), bottomButtonWidth+10, b.screenHeight-bottomBarHeight+4)
}

func (b *BottomBar) Update() error {
	switch entities.Sim.SimulationSpeed {
	case entities.Slow:
		b.bottomButton.Label = ">  "
	case entities.Mid:
		b.bottomButton.Label = ">> "
	default:
		b.bottomButton.Label = ">>>"
	}
	b.bottomButton.Update()

	return nil
}

func NewBottomBar(screenHeight, screenWidth int) *BottomBar {
	return &BottomBar{
		screenHeight: screenHeight,
		screenWidth:  screenWidth,
		bottomButton: Button{
			Label:      ">  ",
			X:          0,
			Y:          screenHeight - 24,
			Width:      bottomButtonWidth,
			Height:     bottomBarHeight,
			Color:      colour.Black,
			HoverColor: colour.Blue,
			OnClick:    entities.Sim.ChangeSimulationSpeed,
		},
	}
}
