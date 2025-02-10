package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/ui/assets"
)

// Sprite represents a drawable entity
type Sprite struct {
	AnimationName string
	FrameIndex    int
	X, Y          float64
}

// SpriteBatch efficiently renders multiple sprites
type SpriteBatch struct {
	Sprites []Sprite
}

// AddSprite queues a sprite for batch rendering
func (b *SpriteBatch) AddSprite(animation string, frameIndex int, x, y float64) {
	b.Sprites = append(b.Sprites, Sprite{AnimationName: animation, FrameIndex: frameIndex, X: x, Y: y})
}

// Draw renders all queued sprites
func (b *SpriteBatch) Draw(screen *ebiten.Image) {
	for _, sprite := range b.Sprites {
		anim, exists := assets.AnimationAssets.Animations[sprite.AnimationName]
		if !exists || len(anim.Frames) == 0 {
			continue // Skip invalid animations
		}

		frame := anim.Frames[sprite.FrameIndex%len(anim.Frames)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sprite.X, sprite.Y)
		screen.DrawImage(frame, op)
	}

	b.Sprites = nil // Clear after drawing
}
