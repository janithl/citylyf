package world

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/janithl/citylyf/internal/entities"
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
)

type WorldRenderer struct {
	playerX, playerY, cameraX, cameraY, zoomFactor float64
	width, height, hoveredTileX, hoveredTileY      int
	startTileX, startTileY                         int
	placingRoad                                    entities.RoadType
}

// Converts grid coordinates to isometric coordinates
func (wr *WorldRenderer) isoTransform(x, y float64) (float64, float64) {
	isoX := (x-y)*float64(tileWidth)/2 + float64(wr.width)/2
	isoY := (x+y)*float64(tileHeight)/4 + float64(wr.height)/4
	return isoX, isoY
}

// Converts screen coordinates to grid coordinates.
func (wr *WorldRenderer) screenToGrid(screenX, screenY float64) (int, int) {
	// Use the same offsets as in isoTransform.
	offsetX := float64(wr.width) / 2
	offsetY := float64(wr.height) / 4

	// Invert the camera/zoom transformation.
	isoX := (screenX-offsetX)/wr.zoomFactor + offsetX + wr.cameraX
	isoY := (screenY-offsetY)/wr.zoomFactor + offsetY + wr.cameraY

	// Invert the isometric transform.
	A := isoX - offsetX
	B := isoY - offsetY

	x := math.Floor(A/float64(tileWidth)+2*B/float64(tileHeight)) - 1
	y := math.Floor(2*B/float64(tileHeight) - A/float64(tileWidth))

	return int(x), int(y)
}

func (wr *WorldRenderer) handleMovement() {
	// Get keyboard input
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		wr.playerX -= moveSpeed
		wr.playerY += moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		wr.playerX += moveSpeed
		wr.playerY -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		wr.playerX -= moveSpeed
		wr.playerY -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		wr.playerX += moveSpeed
		wr.playerY += moveSpeed
	}

	// Smooth camera following
	px, py := wr.isoTransform(wr.playerX, wr.playerY)
	wr.cameraX += (px - float64(wr.width)/2 - wr.cameraX) * 0.1
	wr.cameraY += (py - float64(wr.height)/2 - wr.cameraY) * 0.1
}

func (wr *WorldRenderer) handleZoom() {
	// Mouse wheel zoom
	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		wr.zoomFactor *= 1 + mouseZoomSpeed // Zoom in
	} else if scrollY < 0 {
		wr.zoomFactor *= 1 - mouseZoomSpeed // Zoom out
	}

	// Keyboard zoom
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		wr.zoomFactor *= 1 + kbZoomSpeed // Zoom in
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		wr.zoomFactor *= 1 - kbZoomSpeed // Zoom out
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		wr.zoomFactor = 1 // Reset Zoom
	}

	// Clamp zoom factor between 0.25 and 2
	if wr.zoomFactor < minZoom {
		wr.zoomFactor = minZoom
	} else if wr.zoomFactor > maxZoom {
		wr.zoomFactor = maxZoom
	}
}

func (wr *WorldRenderer) Update() error {
	wr.handleMovement()
	wr.handleZoom()

	// Get mouse position and convert screen coordinates to isometric tile coordinates
	cursorX, cursorY := ebiten.CursorPosition()
	wr.hoveredTileX, wr.hoveredTileY = wr.screenToGrid(float64(cursorX), float64(cursorY))

	// place house
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		entities.Sim.Mutex.Lock()
		entities.Sim.Houses.AddHouse(wr.hoveredTileX, wr.hoveredTileY, 2+rand.Intn(3))
		entities.Sim.Mutex.Unlock()
	}

	// start placing asphalt road
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		wr.placingRoad = entities.Asphalt
		wr.startTileX, wr.startTileY = wr.hoveredTileX, wr.hoveredTileY
	} else if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		wr.placingRoad = entities.Unsealed
		wr.startTileX, wr.startTileY = wr.hoveredTileX, wr.hoveredTileY
	}

	// end placing road
	if wr.placingRoad != "" && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		entities.Sim.Mutex.Lock()
		entities.PlaceRoad(wr.startTileX, wr.startTileY, wr.hoveredTileX, wr.hoveredTileY, wr.placingRoad)
		entities.Sim.Mutex.Unlock()
		wr.placingRoad = ""
	}

	// cancel road placing
	if wr.placingRoad != "" && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		wr.placingRoad = ""
	}

	// toggle roundabout
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		entities.Sim.Mutex.Lock()
		entities.Sim.Geography.ToggleRoundabout(wr.hoveredTileX, wr.hoveredTileY)
		entities.Sim.Mutex.Unlock()
	}

	return nil
}

// Renders the base tile
func (wr *WorldRenderer) renderBaseTiles(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	// Check neighbors (prevent out-of-bounds errors)
	left, right, top, bottom := tiles[x][y].Elevation, tiles[x][y].Elevation, tiles[x][y].Elevation, tiles[x][y].Elevation
	if x > 0 {
		left = tiles[x-1][y].Elevation
	}
	if x < len(tiles)-1 {
		right = tiles[x+1][y].Elevation
	}
	if y > 0 {
		top = tiles[x][y-1].Elevation
	}
	if y < len(tiles[x])-1 {
		bottom = tiles[x][y+1].Elevation
	}

	switch tiles[x][y].Elevation {
	case 8:
		screen.DrawImage(assets.Assets.Sprites["mountain"].Image, op)
	case 7:
		screen.DrawImage(assets.Assets.Sprites["hill"].Image, op)
	case 6:
		if left == 7 && right == 5 {
			screen.DrawImage(assets.Assets.Sprites["slope-x"].Image, op)
		} else if left == 5 && right == 7 {
			screen.DrawImage(assets.Assets.Sprites["slope-x-rev"].Image, op)
		} else if top == 7 && bottom == 5 {
			screen.DrawImage(assets.Assets.Sprites["slope-y"].Image, op)
		} else if top == 5 && bottom == 7 {
			screen.DrawImage(assets.Assets.Sprites["slope-y-rev"].Image, op)
		} else {
			screen.DrawImage(assets.Assets.Sprites["grass"].Image, op)
		}
	case entities.Sim.Geography.SeaLevel:
		screen.DrawImage(assets.Assets.Sprites["sand"].Image, op)
	case 2:
		screen.DrawImage(assets.Assets.Sprites["shallowwater"].Image, op)
	case 1:
		screen.DrawImage(assets.Assets.Sprites["midwater"].Image, op)
	case 0:
		screen.DrawImage(assets.Assets.Sprites["deepwater"].Image, op)
	default:
		screen.DrawImage(assets.Assets.Sprites["grass"].Image, op)
	}

	// Tree Debug
	// if tiles[x][y].Elevation > entities.Sim.Geography.SeaLevel && x%3 == 0 {
	// 	switch y % 4 {
	// 	case 0:
	// 		screen.DrawImage(assets.Assets.Sprites["bush-small"].Image, op)
	// 	case 1:
	// 		screen.DrawImage(assets.Assets.Sprites["bush-large"].Image, op)
	// 	case 2:
	// 		screen.DrawImage(assets.Assets.Sprites["tree-conifer"].Image, op)
	// 	case 3:
	// 		screen.DrawImage(assets.Assets.Sprites["tree-slanted"].Image, op)
	// 	}
	// }
}

// Renders houses
func (wr *WorldRenderer) renderHouses(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	if !tiles[x][y].House { // not a house
		return
	}

	entities.Sim.Mutex.RLock()
	house := entities.Sim.Houses.GetLocationHouse(x, y)
	entities.Sim.Mutex.RUnlock()

	lighting := "dark"
	if house.HouseholdID != 0 {
		lighting = "light"
	}

	if outline, exists := assets.Assets.Sprites[string(house.HouseType)+"-outline-"+lighting]; exists {
		screen.DrawImage(outline.Image, op)
	}

	if houseSprite, exists := assets.Assets.Sprites[string(house.HouseType)+"-"+string(house.RoadDirection)]; exists {
		screen.DrawImage(houseSprite.Image, op)
	}
}

// Renders roads
func (wr *WorldRenderer) renderRoads(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	if !tiles[x][y].Road { // not a road
		return
	}

	roadDirection, roadType := entities.Sim.Geography.IsWithinRoad(x, y)
	roadPrefix := "road-" + string(roadType) + "-"

	// check intersection and draw
	if roadType != "" && tiles[x][y].Intersection != entities.NonIntersection {
		if intersection, exists := assets.Assets.Sprites[roadPrefix+string(tiles[x][y].Intersection)]; exists {
			screen.DrawImage(intersection.Image, op)
		}
	} else {
		// draw correct road
		if roadTile, exists := assets.Assets.Sprites[roadPrefix+string(roadDirection)]; exists {
			screen.DrawImage(roadTile.Image, op)
		}

		// draw correct bridge
		if tiles[x][y].Elevation < entities.Sim.Geography.SeaLevel {
			if bridge, exists := assets.Assets.Sprites["bridge-"+string(roadDirection)]; exists {
				screen.DrawImage(bridge.Image, op)
			}
		}
	}
}

func (wr *WorldRenderer) Draw(screen *ebiten.Image) {
	tiles := entities.Sim.Geography.GetTiles()
	// Use the same offsets as in isoTransform.
	offsetX := float64(wr.width) / 2
	offsetY := float64(wr.height) / 4

	for x := range tiles {
		for y := range tiles[x] {
			isoX, isoY := wr.isoTransform(float64(x), float64(y))

			op := &ebiten.DrawImageOptions{}

			// Apply zoom factor
			op.GeoM.Scale(wr.zoomFactor, wr.zoomFactor)

			// Adjust position using the same offset.
			scaledX := offsetX + (isoX-wr.cameraX-offsetX)*wr.zoomFactor
			scaledY := offsetY + (isoY-wr.cameraY-offsetY)*wr.zoomFactor
			op.GeoM.Translate(scaledX, scaledY)

			wr.renderBaseTiles(screen, op, tiles, x, y)
			wr.renderHouses(screen, op, tiles, x, y)
			wr.renderRoads(screen, op, tiles, x, y)

			// Draw a highlight around the tile under the mouse.
			if wr.hoveredTileX == x && wr.hoveredTileY == y {
				screen.DrawImage(assets.Assets.Sprites["cursorbox-b"].Image, op)
			}

			// Draw a highlight around the tile where the road starts
			if wr.placingRoad != "" && utils.IsWithinRange(wr.startTileX, wr.hoveredTileX, x) && utils.IsWithinRange(wr.startTileY, wr.hoveredTileY, y) {
				screen.DrawImage(assets.Assets.Sprites["cursorbox-r"].Image, op)
			}
		}
	}
}

func NewWorldRenderer(screenWidth, screenHeight int) *WorldRenderer {
	assets.LoadVariableSpritesheet("", "spritesheet-geo.png", "spriteinfo-geo.json")
	assets.LoadVariableSpritesheet("house", "spritesheet-house.png", "spriteinfo-house.json")
	assets.LoadVariableSpritesheet("road", "spritesheet-road.png", "spriteinfo-road.json")

	mapSize := entities.Sim.Geography.Size
	return &WorldRenderer{
		playerX:    float64(mapSize / 3),
		playerY:    float64(mapSize / 3),
		cameraX:    float64(mapSize / 2),
		cameraY:    float64(mapSize / 2),
		zoomFactor: 0.25,
		width:      screenWidth,
		height:     screenHeight,
	}
}
