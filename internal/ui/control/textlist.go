package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/ui/colour"
)

type TextList struct {
	X, Y, Width, Height int
	items               []string
	pages, currentPage  int
	buttons             []*Button
	stepper             *Stepper
	OnClick             func(index int)
}

func (tl *TextList) Update() {
	offset := len(tl.buttons) * (tl.currentPage - 1)
	for i, btn := range tl.buttons {
		j := i + offset
		if j < len(tl.items) {
			btn.Label = tl.items[j]
		} else {
			btn.Label = ""
		}
		btn.Update()
	}
	tl.stepper.SetMaxNumber(tl.pages)
	tl.stepper.Update()
}

func (tl *TextList) Draw(screen *ebiten.Image) {
	for _, btn := range tl.buttons {
		btn.Draw(screen)
	}
	tl.stepper.Draw(screen)
}

func (tl *TextList) SetOffset(x, y int) {
	tl.X = x
	tl.Y = y
	for i, btn := range tl.buttons {
		btn.SetOffset(x, y+i*buttonHeight)
	}
	tl.stepper.SetOffset(tl.getStepperLocation())
}

func (tl *TextList) UpdateItems(items []string) {
	tl.items = items
	tl.pages = (len(items) / len(tl.buttons)) + 1
	if tl.currentPage > tl.pages {
		tl.currentPage = tl.pages
	}
}

// createButtons initializes buttons for the current page
func (tl *TextList) createButtons(count int) {
	tl.buttons = nil // Reset buttons

	for i := 0; i < count; i++ {
		btn := &Button{
			X: tl.X, Y: tl.Y + i*buttonHeight, Width: tl.Width, Height: buttonHeight,
			Label:      "",
			Color:      colour.Transparent,
			HoverColor: colour.SemiBlack,
		}
		btn.OnClick = func() { tl.OnClick(i) }
		tl.buttons = append(tl.buttons, btn)
	}
}

func (tl *TextList) getStepperLocation() (x, y int) {
	stepperX := tl.X + (tl.Width-3*buttonWidth)/2
	stepperY := tl.Y + (tl.Height - buttonHeight)
	return stepperX, stepperY
}

func (tl *TextList) setCurrentPage(page int) {
	if page > 0 && page <= tl.pages {
		tl.currentPage = page
	}
}

// NewTextList creates a list
func NewTextList(x, y, width, height int, items []string) *TextList {
	numberOfButtons := (height / buttonHeight) - 1 // how many buttons can fit into the textlist
	pages := (len(items) / numberOfButtons) + 1    // how many pages of content we have
	tl := &TextList{
		X: x, Y: y, Width: width, Height: height,
		items:       items,
		pages:       pages,
		currentPage: 1,
	}
	stepperX, stepperY := tl.getStepperLocation()
	tl.stepper = NewStepper(stepperX, stepperY, 1, pages, NumberStepper, tl.setCurrentPage)
	tl.createButtons(numberOfButtons)
	return tl
}
