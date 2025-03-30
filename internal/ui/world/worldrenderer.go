package world

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/animation"
	"github.com/janithl/citylyf/internal/ui/assets"
	"github.com/janithl/citylyf/internal/ui/control"
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
	placingUse                         entities.LandUse
	animations                         []*animation.Animation
	tooltip                            *control.Tooltip
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

	// Update the tooltip
	wr.tooltip.Update(ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		wr.tooltip.X, wr.tooltip.Y = cursorX, cursorY
		wr.tooltip.Text = wr.getCursorTileData()
	}

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

	turningPointX, turningPointY := utils.GetTurningPoint(wr.startTile.X, wr.startTile.Y, wr.cursorTile.X, wr.cursorTile.Y)
	for x := range tiles {
		for y := range tiles[x] {
			op := wr.getImageOptions(float64(x), float64(y))

			// draw mountains
			wr.renderMountains(screen, op, tiles, x, y)

			// draw a cursor around the tile under the mouse.
			if x == wr.cursorTile.X && y == wr.cursorTile.Y {
				opCursor := *op
				opCursor.GeoM.Translate(0, wr.elevationToZ(tiles[x][y].Elevation)*wr.zoomFactor) // translate depending on elevation
				screen.DrawImage(assets.Assets.Sprites["ui-cursor"].Image, &opCursor)
			}

			// draw houses and trees last, because they're on the top layer
			wr.renderHouses(screen, op, tiles, x, y)
			wr.renderIndusty(screen, op, tiles, x, y)

			op.GeoM.Translate(0, wr.elevationToZ(tiles[x][y].Elevation)*wr.zoomFactor) // translate depending on elevation

			// draw zones
			if tiles[x][y].LandUse != entities.NoUse && tiles[x][y].LandStatus != entities.DevelopedStatus {
				if zone, exists := assets.Assets.Sprites["ui-zone-"+string(tiles[x][y].LandUse)]; exists {
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
			landUseHighlight := wr.placingUse != entities.NoUse &&
				utils.IsWithinRange(wr.startTile.X, wr.cursorTile.X, x) && utils.IsWithinRange(wr.startTile.Y, wr.cursorTile.Y, y)
			roadHighlight := wr.placingRoad != entities.NoRoad &&
				((utils.IsWithinRange(wr.startTile.X, turningPointX, x) && utils.IsWithinRange(wr.startTile.Y, turningPointY, y)) ||
					(utils.IsWithinRange(wr.cursorTile.X, turningPointX, x) && utils.IsWithinRange(wr.cursorTile.Y, turningPointY, y)))

			if landUseHighlight || roadHighlight {
				if tiles[x][y].IsBuildable() {
					screen.DrawImage(assets.Assets.Sprites["ui-highlight"].Image, op)
				} else {
					screen.DrawImage(assets.Assets.Sprites["ui-highlight-danger"].Image, op)
				}
			}
		}
	}
	wr.tooltip.Draw(screen)
}

func NewWorldRenderer(screenWidth, screenHeight int) *WorldRenderer {
	assets.LoadVariableSpritesheet("", "spritesheet-geo.png", "spriteinfo-geo.json")
	assets.LoadVariableSpritesheet("house", "spritesheet-house.png", "spriteinfo-house.json")
	assets.LoadVariableSpritesheet("industry", "spritesheet-industry.png", "spriteinfo-industry.json")
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
		tooltip: &control.Tooltip{
			Height:  72,
			Width:   210,
			Padding: 4,
			Margin:  20,
			Text:    "Tooltip",
		},
	}
}
