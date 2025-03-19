package world

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/animation"
	"github.com/janithl/citylyf/internal/ui/assets"
	"github.com/janithl/citylyf/internal/utils"
)

const (
	tileWidth      = 64
	tileHeight     = 64
	moveSpeed      = 0.2
	mouseZoomSpeed = 0.05
	kbZoomSpeed    = 0.035
	minZoom        = 0.25
	maxZoom        = 2
	totalAnims     = 256
)

var animatedHumans = []string{"teal", "green", "orange", "pink"}

type WorldRenderer struct {
	playerX, playerY, offsetX, offsetY float64
	cameraX, cameraY, zoomFactor       float64
	width, height, frameCounter        int
	cursorTile, startTile              entities.Point
	placingRoad                        entities.RoadType
	placingZone                        entities.Zone
	animations                         []*animation.Animation
}

func (wr *WorldRenderer) Update(mapRegenMode bool) error {
	if mapRegenMode {
		return nil
	}

	wr.frameCounter++
	if wr.frameCounter >= 300 { // update every 5 seconds
		wr.frameCounter = 0

		// assign animations
		entities.Sim.Mutex.RLock()
		wr.assignAnimations()
		entities.Sim.Mutex.RUnlock()
	}

	for i := range wr.animations {
		wr.animations[i].Update()
	}

	wr.handleMovement()
	wr.handleZoom()

	// Get mouse position and convert screen coordinates to isometric tile coordinates
	cursorX, cursorY := ebiten.CursorPosition()
	wr.cursorTile = wr.screenToGrid(float64(cursorX), float64(cursorY))
	wr.getUserInput()

	return nil
}

func (wr *WorldRenderer) Draw(screen *ebiten.Image) {
	tiles := entities.Sim.Geography.GetTiles()
	for x := range tiles {
		for y := range tiles[x] {
			op := wr.getImageOptions(float64(x), float64(y))
			wr.renderBaseTiles(screen, op, tiles, x, y)
			wr.renderRoads(screen, op, tiles, x, y)
		}
	} // finish rendering base tiles first to prevent overlaps with everything else

	for x := range tiles {
		for y := range tiles[x] {
			op := wr.getImageOptions(float64(x), float64(y))

			// draw a cursor around the tile under the mouse.
			if x == wr.cursorTile.X && y == wr.cursorTile.Y {
				opCursor := wr.getImageOptions(float64(x), float64(y))
				opCursor.GeoM.Translate(0, wr.elevationToZ(tiles[x][y].Elevation)*wr.zoomFactor) // translate depending on elevation
				screen.DrawImage(assets.Assets.Sprites["ui-cursor"].Image, opCursor)
			}

			// draw mountains, houses and trees last, because they're on the top layer
			wr.renderMountains(screen, op, tiles, x, y)
			wr.renderHouses(screen, op, tiles, x, y)
			wr.renderRetail(screen, op, tiles, x, y)

			op.GeoM.Translate(0, wr.elevationToZ(tiles[x][y].Elevation)*wr.zoomFactor) // translate depending on elevation
			// draw zones
			if tiles[x][y].Zone != entities.NoZone {
				if zone, exists := assets.Assets.Sprites["ui-"+string(tiles[x][y].Zone)]; exists {
					screen.DrawImage(zone.Image, op)
				}
			}

			// draw tile animations
			for _, anim := range wr.animations {
				if animX, animY := anim.Coordinates(); x == animX && y == animY {
					anim.Draw(screen, wr.getImageOptions)
				}
			}

			// draw a highlight around the tile where the road starts
			if (wr.placingRoad != entities.NoRoad || wr.placingZone != entities.NoZone) &&
				utils.IsWithinRange(wr.startTile.X, wr.cursorTile.X, x) && utils.IsWithinRange(wr.startTile.Y, wr.cursorTile.Y, y) {
				screen.DrawImage(assets.Assets.Sprites["ui-highlight"].Image, op)
			}
		}
	}

}

func NewWorldRenderer(screenWidth, screenHeight int) *WorldRenderer {
	assets.LoadVariableSpritesheet("", "spritesheet-geo.png", "spriteinfo-geo.json")
	assets.LoadVariableSpritesheet("house", "spritesheet-house.png", "spriteinfo-house.json")
	assets.LoadVariableSpritesheet("retail", "spritesheet-retail.png", "spriteinfo-retail.json")
	assets.LoadVariableSpritesheet("road", "spritesheet-road.png", "spriteinfo-road.json")
	assets.LoadVariableSpritesheet("ui", "spritesheet-ui.png", "spriteinfo-ui.json")

	animations := make([]*animation.Animation, totalAnims) // support up to n walking animations at a time
	for i := range animations {
		animations[i] = animation.NewAnimation(animatedHumans[rand.IntN(len(animatedHumans))], 0, 0)
	}

	mapSize := entities.Sim.Geography.Size
	return &WorldRenderer{
		playerX:      float64(mapSize / 3),
		playerY:      float64(mapSize / 3),
		cameraX:      float64(mapSize / 2),
		cameraY:      float64(mapSize / 2),
		zoomFactor:   0.25,
		width:        screenWidth,
		height:       screenHeight,
		offsetX:      float64(screenWidth) / 2,
		offsetY:      float64(screenHeight) / 4,
		animations:   animations,
		frameCounter: 300,
	}
}
