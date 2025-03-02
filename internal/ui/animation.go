package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/ui/assets"
)

type Animation struct {
	X, Y, SpeedX, SpeedY float64
	batch                SpriteBatch
	frameIndex           int
	frameCount           int
}

func (a *Animation) Update() error {
	a.frameCount++
	if a.frameCount%8 == 0 { // Change frame every 8 ticks
		a.frameIndex++
	}

	// Move sprite
	a.X += a.SpeedX
	a.Y += a.SpeedY

	return nil
}

func (a *Animation) Draw(screen *ebiten.Image) {
	// Queue animated character rendering based on direction of walk
	switch {
	case a.SpeedX > 0 && a.SpeedY > 0:
		a.batch.AddSprite("walk_front_2", a.frameIndex, a.X, a.Y)
	case a.SpeedX < 0 && a.SpeedY > 0:
		a.batch.AddSprite("walk_front_3", a.frameIndex, a.X, a.Y)
	case a.SpeedX < 0 && a.SpeedY == 0:
		a.batch.AddSprite("walk_side_2", a.frameIndex, a.X, a.Y)
	case a.SpeedX > 0 && a.SpeedY == 0:
		a.batch.AddSprite("walk_side_1", a.frameIndex, a.X, a.Y)
	case a.SpeedX > 0 && a.SpeedY < 0:
		a.batch.AddSprite("walk_back_1", a.frameIndex, a.X, a.Y)
	case a.SpeedX == 0 && a.SpeedY < 0:
		a.batch.AddSprite("walk_back_2", a.frameIndex, a.X, a.Y)
	case a.SpeedX < 0 && a.SpeedY < 0:
		a.batch.AddSprite("walk_back_3", a.frameIndex, a.X, a.Y)
	case a.SpeedX == 0 && a.SpeedY == 0:
		a.batch.AddSprite("walk_front_1", 0, a.X, a.Y)
	default:
		a.batch.AddSprite("walk_front_1", a.frameIndex, a.X, a.Y)
	}

	// Execute batch render
	a.batch.Draw(screen)
}

func NewAnimation(x, y, speedX, speedY float64) *Animation {
	animations := map[string]int{
		"walk_front_1": 0,
		"walk_front_2": 1,
		"walk_side_1":  2,
		"walk_back_1":  3,
		"walk_back_2":  4,
		"walk_back_3":  5,
		"walk_side_2":  6,
		"walk_front_3": 7,
	}

	// shout out to https://bossnelnel.itch.io/8-direction-top-down-character-sprites for the amazing sprites
	assets.LoadAnimationSpritesheet("human-green.png", 23, 36, 9, 8, animations)
	return &Animation{
		X:      x,
		Y:      y,
		SpeedX: speedX,
		SpeedY: speedY,
	}
}
