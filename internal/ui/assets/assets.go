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
	Animations map[string]Animation // Stores animations by name
	Sprites    map[string]Sprite
	Names      map[string][]string
}

// Global instance
var Assets *AssetManager

func init() {
	// Initialize asset manager
	Assets = &AssetManager{
		Sprites:    make(map[string]Sprite),
		Animations: make(map[string]Animation),
		Names:      make(map[string][]string),
	}
	LoadNames()
}

// LoadSpritesheet loads a multi-line sprite sheet
func LoadAnimationSpritesheet(prefix, imagePath string, frameWidth, frameHeight, columns, rows int, animations map[string]int) {
	img, _, err := ebitenutil.NewImageFromFileSystem(assetsFolder, imagePath)
	if err != nil {
		log.Fatal(err)
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
		Assets.Animations[prefix+"_"+name] = Animation{Frames: frames}
	}
}

// LoadVariableSpritesheet loads a spritesheet and extracts sprites using JSON definitions
func LoadVariableSpritesheet(prefix, imagePath, jsonPath string) {
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

	// add a hyphen
	if prefix != "" {
		prefix += "-"
	}

	// Extract sprites based on JSON data
	for name, rect := range spriteMap {
		sprite := img.SubImage(image.Rect(rect.X, rect.Y, rect.X+rect.Width, rect.Y+rect.Height)).(*ebiten.Image)
		Assets.Sprites[prefix+name] = Sprite{
			Image: sprite,
			X:     rect.X,
			Y:     rect.Y,
		}
	}
}

// LoadNames loads the names.json
func LoadNames() {
	// Load JSON file
	data, err := fs.ReadFile(assetsFolder, "names.json")
	if err != nil {
		log.Fatal("Failed to load names JSON:", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &Assets.Names); err != nil {
		log.Fatal("Failed to parse sprite JSON:", err)
	}
}
