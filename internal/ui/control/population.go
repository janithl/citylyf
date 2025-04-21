package control

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type PopulationPyramid struct {
	X, Y, Width, Height, frameCounter, maxPopPerGroup int
	ageGroups                                         map[int]entities.AgeGroup
}

const BarGraphPadding = 4

func (pp *PopulationPyramid) Draw(screen *ebiten.Image) {
	barHeight := float32(pp.Height) / float32(len(pp.ageGroups)) // Divide height by number of age groups
	maxWidth := float32(pp.Width) / 2                            // Half width for each side

	for i, group := range pp.ageGroups {
		y := float32(pp.Y) + float32(pp.Height) - barHeight - float32(i/entities.AgeGroupSize)*barHeight
		maleWidth := maxWidth * (float32(group.Male) / float32(pp.maxPopPerGroup))
		femaleWidth := maxWidth * (float32(group.Female) / float32(pp.maxPopPerGroup))
		otherWidth := maxWidth * (float32(group.Other) / float32(pp.maxPopPerGroup))

		// Draw bars
		vector.DrawFilledRect(screen, float32(pp.X)+maxWidth-maleWidth, y, maleWidth, barHeight-BarGraphPadding, colour.DarkCyan, false)
		vector.DrawFilledRect(screen, float32(pp.X)+maxWidth, y, femaleWidth, barHeight-BarGraphPadding, colour.DarkMagenta, false)
		vector.DrawFilledRect(screen, float32(pp.X)+maxWidth-otherWidth/2, y, otherWidth, barHeight-BarGraphPadding, colour.Gray, false)

		// Label age group
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%3d - %-3d", i, i+entities.AgeGroupSize), int(pp.X)+int(pp.Width)/2-28, int(y))
	}
}

func (pp *PopulationPyramid) Update() {
	pp.frameCounter++
	if pp.frameCounter >= 60 { // update every second
		pp.frameCounter = 0

		entities.Sim.Mutex.RLock()
		pp.ageGroups = entities.Sim.People.AgeGroups
		pp.maxPopPerGroup = entities.Sim.People.Population()
		if len(pp.ageGroups) > 0 && entities.Sim.People.Population() > 20 { // bigger populations are easier to predict
			pp.maxPopPerGroup = 3 * entities.Sim.People.Population() / len(pp.ageGroups)
		}
		entities.Sim.Mutex.RUnlock()
	}
}

func (pp *PopulationPyramid) SetOffset(x, y int) {
	pp.X = x
	pp.Y = y
}
