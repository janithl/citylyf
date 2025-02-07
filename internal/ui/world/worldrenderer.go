package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/ui/assets"
)

const (
	tileWidth  = 64
	tileHeight = 32
	moveSpeed  = 0.1
)

type Tile int

const (
	Grass Tile = 0
	Road  Tile = 1
)

type WorldRenderer struct {
	playerX, playerY, cameraX, cameraY float64
	width, height                      int
	tiles                              [][]Tile
	grassImage                         *ebiten.Image
	roadImage                          *ebiten.Image
}

// Converts grid coordinates to isometric coordinates
func (wr *WorldRenderer) isoTransform(x, y float64) (float64, float64) {
	isoX := (x-y)*float64(tileWidth)/2 + float64(wr.width)/2
	isoY := (x+y)*float64(tileHeight)/2 + float64(wr.height)/4
	return isoX, isoY
}

func (wr *WorldRenderer) Update() error {
	// Handle movement
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

	return nil
}

func (wr *WorldRenderer) Draw(screen *ebiten.Image) {
	// Draw isometric tiles
	for x := range wr.tiles {
		for y := range wr.tiles[x] {
			isoX, isoY := wr.isoTransform(float64(x), float64(y))

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(isoX-wr.cameraX, isoY-wr.cameraY)

			if wr.tiles[x][y] == Grass {
				screen.DrawImage(wr.grassImage, op)
			} else {
				screen.DrawImage(wr.roadImage, op)
			}
		}
	}
}

func NewWorldRenderer(screenWidth, screenHeight, mapWidth, mapHeight int) *WorldRenderer {
	tiles := make([][]Tile, mapWidth)
	for x := range tiles {
		tiles[x] = make([]Tile, mapHeight)
	}

	tiles[3][5] = Road
	tiles[3][6] = Road
	tiles[3][7] = Road

	return &WorldRenderer{
		playerX:    float64(mapWidth / 2),
		playerY:    float64(mapHeight / 2),
		cameraX:    float64(mapWidth / 2),
		cameraY:    float64(mapHeight / 2),
		width:      screenWidth,
		height:     screenHeight,
		grassImage: assets.LoadImage("internal/ui/assets/grass.png"),
		roadImage:  assets.LoadImage("internal/ui/assets/road.png"),
		tiles:      tiles,
	}
}
