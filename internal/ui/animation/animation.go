package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/assets"
)

type Animation struct {
	x, y, speedX, speedY   float64
	path                   []*entities.Point
	batch                  *SpriteBatch
	prefix                 string
	frameIndex, frameCount int
}

func (a *Animation) Update() error {
	a.frameCount++
	if a.frameCount%8 == 0 { // Change frame every 8 ticks
		a.frameIndex++
	}

	if len(a.path) > 1 && a.path[1].X-int(a.x) < 3 && a.path[1].Y-int(a.y) < 3 {
		a.path = a.path[1:]
		a.CalculateSpeed()
	}

	// Move sprite
	a.x += a.speedX
	a.y += a.speedY

	return nil
}

func (a *Animation) Draw(screen *ebiten.Image, getImageOptions func(float64, float64) *ebiten.DrawImageOptions) {
	if len(a.path) <= 1 {
		return
	}

	// Queue animated character rendering based on direction of walk
	switch {
	case a.speedX > 0 && a.speedY > 0:
		a.batch.AddSprite(a.prefix+"_walk_front_2", a.frameIndex, a.x, a.y)
	case a.speedX < 0 && a.speedY > 0:
		a.batch.AddSprite(a.prefix+"_walk_front_3", a.frameIndex, a.x, a.y)
	case a.speedX < 0 && a.speedY == 0:
		a.batch.AddSprite(a.prefix+"_walk_side_2", a.frameIndex, a.x, a.y)
	case a.speedX > 0 && a.speedY == 0:
		a.batch.AddSprite(a.prefix+"_walk_side_1", a.frameIndex, a.x, a.y)
	case a.speedX > 0 && a.speedY < 0:
		a.batch.AddSprite(a.prefix+"_walk_back_1", a.frameIndex, a.x, a.y)
	case a.speedX == 0 && a.speedY < 0:
		a.batch.AddSprite(a.prefix+"_walk_back_2", a.frameIndex, a.x, a.y)
	case a.speedX < 0 && a.speedY < 0:
		a.batch.AddSprite(a.prefix+"_walk_back_3", a.frameIndex, a.x, a.y)
	case a.speedX == 0 && a.speedY == 0:
		a.batch.AddSprite(a.prefix+"_walk_front_1", 0, a.x, a.y)
	default:
		a.batch.AddSprite(a.prefix+"_walk_front_1", a.frameIndex, a.x, a.y)
	}

	// Execute batch render
	a.batch.Draw(screen, getImageOptions)
}

func (a *Animation) CalculateSpeed() {
	if len(a.path) <= 1 {
		a.speedX, a.speedY = 0, 0
		return
	}

	if a.path[1].X-a.path[0].X > 0 && a.path[1].Y-a.path[0].Y > 0 {
		a.speedX, a.speedY = 0.035, 0.035
	} else if a.path[1].X-a.path[0].X > 0 {
		a.speedX = 0.05
	} else if a.path[1].Y-a.path[0].Y > 0 {
		a.speedY = 0.05
	} else {
		a.speedX, a.speedY = 0, 0
	}
}

func NewAnimation(prefix string, x, y float64, path []*entities.Point) *Animation {
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
	assets.LoadAnimationSpritesheet(prefix, "human-"+prefix+".png", 23, 36, 9, 8, animations)

	anim := &Animation{
		x:      x,
		y:      y,
		path:   path,
		batch:  &SpriteBatch{},
		prefix: prefix,
	}
	anim.CalculateSpeed()
	return anim
}
