package colour

import "image/color"

// Basic Colors
var (
	White         = color.RGBA{255, 255, 255, 255}
	Black         = color.RGBA{0, 0, 0, 255}
	Gray          = color.RGBA{128, 128, 128, 255}
	LightGray     = color.RGBA{192, 192, 192, 255}
	DarkGray      = color.RGBA{64, 64, 64, 255}
	Red           = color.RGBA{255, 0, 0, 255}
	DarkRed       = color.RGBA{128, 0, 0, 255}
	Green         = color.RGBA{0, 255, 0, 255}
	DarkGreen     = color.RGBA{0, 128, 0, 255}
	Blue          = color.RGBA{0, 0, 255, 255}
	DarkBlue      = color.RGBA{0, 0, 128, 255}
	Yellow        = color.RGBA{255, 255, 0, 255}
	Cyan          = color.RGBA{0, 255, 255, 255}
	DarkCyan      = color.RGBA{0, 192, 192, 255}
	Magenta       = color.RGBA{255, 0, 255, 255}
	DarkMagenta   = color.RGBA{192, 0, 192, 255}
	Transparent   = color.RGBA{0, 0, 0, 0}
	SemiBlack     = color.RGBA{0, 0, 0, 128}
	DarkSemiBlack = color.RGBA{0, 0, 0, 192}
)
