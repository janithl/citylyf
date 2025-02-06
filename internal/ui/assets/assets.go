package assets

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Animation stores frames for a specific action
type Animation struct {
	Frames []*ebiten.Image
}

// AssetManager manages all game assets
type AssetManager struct {
	SpriteSheet *ebiten.Image
	Animations  map[string]Animation // Stores animations by name
}

// Global instance
var Assets *AssetManager

// LoadSpritesheet loads a multi-line sprite sheet
func LoadSpritesheet(path string, frameWidth, frameHeight, columns, rows int, animations map[string]int) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	Assets = &AssetManager{
		SpriteSheet: img,
		Animations:  make(map[string]Animation),
	}

	// Extract animations based on row mappings
	for name, row := range animations {
		var frames []*ebiten.Image
		for col := 0; col < columns; col++ {
			frame := img.SubImage(image.Rect(
				col*frameWidth, row*frameHeight, (col+1)*frameWidth, (row+1)*frameHeight,
			)).(*ebiten.Image)
			frames = append(frames, frame)
		}
		Assets.Animations[name] = Animation{Frames: frames}
	}
}
