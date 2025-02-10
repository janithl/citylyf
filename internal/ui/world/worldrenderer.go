package world

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/assets"
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

	return nil
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

			switch tiles[x][y].Elevation {
			case 8:
				screen.DrawImage(assets.Assets.Sprites["mountain"].Image, op)
			case 7, 6:
				screen.DrawImage(assets.Assets.Sprites["hill"].Image, op)
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

			// Draw roads if necessary
			roadDirection := entities.Sim.Geography.IsWithinRoad(x, y)
			if roadDirection == entities.DirX {
				screen.DrawImage(assets.Assets.Sprites["road-x"].Image, op)
			} else if roadDirection == entities.DirY {
				screen.DrawImage(assets.Assets.Sprites["road-y"].Image, op)
			}

			// Draw a highlight around the tile under the mouse.
			if wr.hoveredTileX == x && wr.hoveredTileY == y {
				screen.DrawImage(assets.Assets.Sprites["cursorbox"].Image, op)
			}
		}
	}
}

func NewWorldRenderer(screenWidth, screenHeight int) *WorldRenderer {
	assets.LoadVariableSpritesheet("internal/ui/assets/geo-spritesheet.png", "internal/ui/assets/sprites.json")

	mapSize := entities.Sim.Geography.Size
	return &WorldRenderer{
		playerX:    float64(mapSize / 2),
		playerY:    float64(mapSize / 2),
		cameraX:    float64(mapSize / 2),
		cameraY:    float64(mapSize / 2),
		zoomFactor: 1,
		width:      screenWidth,
		height:     screenHeight,
	}
}
