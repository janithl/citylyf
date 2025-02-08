package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/assets"
)

const (
	tileWidth  = 64
	tileHeight = 64
	moveSpeed  = 0.1
)

type WorldRenderer struct {
	playerX, playerY, cameraX, cameraY float64
	width, height                      int
}

// Converts grid coordinates to isometric coordinates
func (wr *WorldRenderer) isoTransform(x, y float64) (float64, float64) {
	isoX := (x-y)*float64(tileWidth)/2 + float64(wr.width)/2
	isoY := (x+y)*float64(tileHeight)/4 + float64(wr.height)/4
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
	tiles := entities.Sim.Geography.GetTiles()
	for x := range tiles {
		for y := range tiles[x] {
			isoX, isoY := wr.isoTransform(float64(x), float64(y))

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(isoX-wr.cameraX, isoY-wr.cameraY)

			switch tiles[x][y].Elevation {
			case 7:
				screen.DrawImage(assets.Assets.Sprites["mountain"].Image, op)
			case 6:
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
		}
	}
}

func NewWorldRenderer(screenWidth, screenHeight int) *WorldRenderer {
	assets.LoadVariableSpritesheet("internal/ui/assets/geo-spritesheet.png", "internal/ui/assets/sprites.json")

	mapSize := entities.Sim.Geography.Size
	return &WorldRenderer{
		playerX: float64(mapSize / 2),
		playerY: float64(mapSize / 2),
		cameraX: float64(mapSize / 2),
		cameraY: float64(mapSize / 2),
		width:   screenWidth,
		height:  screenHeight,
	}
}
