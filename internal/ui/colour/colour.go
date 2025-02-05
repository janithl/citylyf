package colour

import "image/color"

// Basic Colors
var (
	White       = color.RGBA{255, 255, 255, 255}
	Black       = color.RGBA{0, 0, 0, 255}
	Gray        = color.RGBA{128, 128, 128, 255}
	LightGray   = color.RGBA{196, 196, 196, 255}
	DarkGray    = color.RGBA{64, 64, 64, 255}
	Red         = color.RGBA{255, 0, 0, 255}
	Green       = color.RGBA{0, 255, 0, 255}
	Blue        = color.RGBA{0, 0, 255, 255}
	Yellow      = color.RGBA{255, 255, 0, 255}
	Cyan        = color.RGBA{0, 255, 255, 255}
	Magenta     = color.RGBA{255, 0, 255, 255}
	Transparent = color.RGBA{0, 0, 0, 0}
	SemiBlack   = color.RGBA{0, 0, 0, 128}
)
