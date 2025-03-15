package animation

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/assets"
)

type Animation struct {
	x, y, speedX, speedY            float64
	path                            []*entities.Point
	batch                           *SpriteBatch
	prefix                          string
	frameIndex, frameCounter, delay int
	finished                        bool
}

func (a *Animation) Update() error {
	if len(a.path) <= 1 {
		a.finished = true
		return nil
	}

	entities.Sim.Mutex.RLock()
	simSpeed := float64(entities.Sim.SimulationSpeed)
	entities.Sim.Mutex.RUnlock()

	if simSpeed > 0 { // Pause walking if sim is paused
		a.frameCounter += int(math.Sqrt(1600 / simSpeed))
	}

	if a.delay > 0 { // don't render if we are still delayed
		a.delay--
		return nil
	}

	if a.frameCounter > 8 { // Change frame every 8 ticks
		a.frameIndex++
		a.frameCounter = 0
	}

	// Move to the next path point if close enough
	dx := a.path[1].X - int(a.x)
	dy := a.path[1].Y - int(a.y)
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	if distance < 0.1 { // Threshold for reaching the point
		a.path = a.path[1:]
		a.CalculateSpeed(a.delay)
	}

	// Move sprite along the speed vector based on sim speed
	if simSpeed > 0 {
		a.x += a.speedX / simSpeed
		a.y += a.speedY / simSpeed
	}

	return nil
}

func (a *Animation) Draw(screen *ebiten.Image, getImageOptions func(float64, float64) *ebiten.DrawImageOptions) {
	if a.finished || a.delay > 0 {
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

func (a *Animation) CalculateSpeed(delay int) {
	if len(a.path) <= 1 {
		a.speedX, a.speedY = 0, 0
		return
	}

	a.delay = delay

	// Get direction vector
	dx := float64(a.path[1].X - a.path[0].X)
	dy := float64(a.path[1].Y - a.path[0].Y)
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance == 0 {
		a.speedX, a.speedY = 0, 0
		return
	}

	// Normalize direction and apply speed
	speed := 16.0 // Base speed value
	a.speedX = (dx / distance) * speed
	a.speedY = (dy / distance) * speed
}

func (a *Animation) IsFinished() bool {
	return a.finished
}

func (a *Animation) Coordinates() (int, int) {
	return int(math.Round(a.x)), int(math.Round(a.y))
}

func (a *Animation) SetPath(path []*entities.Point) {
	if len(path) > 0 {
		a.path = path
		a.finished = false
		a.x, a.y = float64(path[0].X), float64(path[0].Y)
	}
}

func NewAnimation(prefix string, x, y float64) *Animation {
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
	assets.LoadAnimationSpritesheet(prefix, "human-"+prefix+".png", 64, 64, 9, 8, animations)

	anim := &Animation{
		x:        x,
		y:        y,
		batch:    &SpriteBatch{},
		prefix:   prefix,
		finished: true,
	}
	return anim
}
