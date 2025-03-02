package assets

import (
	"embed"
	"encoding/json"
	"image"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed *.png *.json
var assetsFolder embed.FS

// Sprite represents an extracted image from the spritesheet
type Sprite struct {
	Image               *ebiten.Image
	X, Y, Width, Height int
}

// Animation stores frames for a specific action
type Animation struct {
	Frames []*ebiten.Image
}

// AssetManager manages all game assets
type AssetManager struct {
	SpriteSheet *ebiten.Image
	Animations  map[string]Animation // Stores animations by name
	Sprites     map[string]Sprite
}

// Global instance
var AnimationAssets *AssetManager
var Assets *AssetManager

// LoadSpritesheet loads a multi-line sprite sheet
func LoadAnimationSpritesheet(path string, frameWidth, frameHeight, columns, rows int, animations map[string]int) {
	img, _, err := ebitenutil.NewImageFromFileSystem(assetsFolder, path)
	if err != nil {
		log.Fatal(err)
	}

	AnimationAssets = &AssetManager{
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
		AnimationAssets.Animations[name] = Animation{Frames: frames}
	}
}

// Load Single Image
func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

// LoadVariableSpritesheet loads a spritesheet and extracts sprites using JSON definitions
func LoadVariableSpritesheet(imagePath, jsonPath string) {
	// Load JSON file
	data, err := fs.ReadFile(assetsFolder, jsonPath)
	if err != nil {
		log.Fatal("Failed to load sprite JSON:", err)
	}

	// Parse JSON into a map
	var spriteMap map[string]struct {
		X, Y, Width, Height int
	}

	if err := json.Unmarshal(data, &spriteMap); err != nil {
		log.Fatal("Failed to parse sprite JSON:", err)
	}

	// Load the spritesheet image
	img, _, err := ebitenutil.NewImageFromFileSystem(assetsFolder, imagePath)
	if err != nil {
		log.Fatal("Failed to load spritesheet:", err)
	}

	// Initialize asset manager
	Assets = &AssetManager{
		SpriteSheet: img,
		Sprites:     make(map[string]Sprite),
	}

	// Extract sprites based on JSON data
	for name, rect := range spriteMap {
		sprite := img.SubImage(image.Rect(rect.X, rect.Y, rect.X+rect.Width, rect.Y+rect.Height)).(*ebiten.Image)
		Assets.Sprites[name] = Sprite{
			Image: sprite,
			X:     rect.X,
			Y:     rect.Y,
		}
	}
}
