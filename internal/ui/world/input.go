package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/janithl/citylyf/internal/entities"
)

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

func (wr *WorldRenderer) getUserInput() {
	// start placing zone
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		wr.placingRoad = entities.NoRoad
		wr.placingUse = entities.ResidentialUse
		wr.startTile = entities.Point{X: wr.cursorTile.X, Y: wr.cursorTile.Y}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyY) {
		wr.placingRoad = entities.NoRoad
		wr.placingUse = entities.RetailUse
		wr.startTile = entities.Point{X: wr.cursorTile.X, Y: wr.cursorTile.Y}
	}

	// start placing asphalt road
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		wr.placingUse = entities.NoUse
		wr.placingRoad = entities.Asphalt
		wr.startTile = entities.Point{X: wr.cursorTile.X, Y: wr.cursorTile.Y}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		wr.placingUse = entities.NoUse
		wr.placingRoad = entities.Unsealed
		wr.startTile = entities.Point{X: wr.cursorTile.X, Y: wr.cursorTile.Y}
	}

	// end placing road/zone
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if wr.placingRoad != entities.NoRoad {
			entities.Sim.Mutex.Lock()
			entities.PlaceRoad(wr.startTile, wr.cursorTile, wr.placingRoad)
			entities.Sim.Mutex.Unlock()
			wr.placingRoad = entities.NoRoad
		} else if wr.placingUse != entities.NoUse {
			entities.Sim.Mutex.Lock()
			entities.Sim.Geography.PlaceLandUse(wr.startTile, wr.cursorTile, wr.placingUse)
			entities.Sim.Mutex.Unlock()
			wr.placingUse = entities.NoUse
		}
	}

	// cancel road/zone placing
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		wr.placingRoad = entities.NoRoad
		wr.placingUse = entities.NoUse
	}

	// toggle roundabout
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		entities.Sim.Mutex.Lock()
		entities.Sim.Geography.ToggleRoundabout(wr.cursorTile.X, wr.cursorTile.Y)
		entities.Sim.Mutex.Unlock()
	}
}
