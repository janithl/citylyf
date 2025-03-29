package world

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/utils"
)

// Converts grid coordinates to isometric coordinates
func (wr *WorldRenderer) isoTransform(x, y float64) (float64, float64) {
	isoX := (x-y)*float64(tileWidth)/2 + float64(wr.width)/2
	isoY := (x+y)*float64(tileHeight)/4 + float64(wr.height)/4
	return isoX, isoY
}

// Converts screen coordinates to grid coordinates.
func (wr *WorldRenderer) screenToGrid(screenX, screenY float64) entities.Point {
	// Invert the camera/zoom transformation.
	isoX := (screenX-wr.offsetX)/wr.zoomFactor + wr.offsetX + wr.cameraX
	isoY := (screenY-wr.offsetY)/wr.zoomFactor + wr.offsetY + wr.cameraY

	// Invert the isometric transform.
	A := isoX - wr.offsetX
	B := isoY - wr.offsetY

	x := math.Floor(A/float64(tileWidth)+2*B/float64(tileHeight)) - 1
	y := math.Floor(2*B/float64(tileHeight) - A/float64(tileWidth))

	return entities.Point{X: int(x), Y: int(y)}
}

// converts elevation to screen position changes
func (wr *WorldRenderer) elevationToZ(elevation int) float64 {
	switch {
	case elevation < entities.Sim.Geography.SeaLevel:
		return 0
	case elevation >= entities.Sim.Geography.HillLevel:
		return -16
	default:
		return -8
	}
}

// returns ebiten image options for a given (x,y) coordinate
func (wr *WorldRenderer) getImageOptions(x, y float64) *ebiten.DrawImageOptions {
	isoX, isoY := wr.isoTransform(x, y)

	op := &ebiten.DrawImageOptions{}

	// Apply zoom factor
	op.GeoM.Scale(wr.zoomFactor, wr.zoomFactor)

	// Adjust position using the same offset.
	scaledX := wr.offsetX + (isoX-wr.cameraX-wr.offsetX)*wr.zoomFactor
	scaledY := wr.offsetY + (isoY-wr.cameraY-wr.offsetY)*wr.zoomFactor
	op.GeoM.Translate(scaledX, scaledY)

	return op
}

// getCursorTileData returns current cursor tile data
func (wr *WorldRenderer) getCursorTileData() string {
	entities.Sim.Mutex.RLock()
	tiles := entities.Sim.Geography.GetTiles()
	entities.Sim.Mutex.RUnlock()

	if wr.cursorTile.X >= 0 && wr.cursorTile.X < len(tiles) && wr.cursorTile.Y >= 0 && wr.cursorTile.Y < len(tiles) {
		tile := tiles[wr.cursorTile.X][wr.cursorTile.Y]
		built := ""
		if tile.IsBuilt() {
			output := ""
			entities.Sim.Mutex.RLock()
			if tile.LandUse == entities.ResidentialUse {
				if house := entities.Sim.Houses.GetLocationHouse(wr.cursorTile.X, wr.cursorTile.Y); house != nil {
					output = fmt.Sprintf("#%d: %d Bedroom House\nRent: $%d/month", house.ID, house.Bedrooms, house.MonthlyRent)
				}
			} else if tile.LandUse == entities.RetailUse || tile.LandUse == entities.AgricultureUse {
				if company := entities.Sim.Companies.GetLocationCompany(wr.cursorTile.X, wr.cursorTile.Y); company != nil {
					output = fmt.Sprintf("#%d: %s\n%d Employees / %d Openings\nProfit/Loss: %s\nRevenue: %s", company.ID, company.Name,
						company.GetNumberOfEmployees(), company.GetNumberOfJobOpenings(),
						utils.FormatCurrency(company.LastProfit, "$"),
						utils.FormatCurrency(company.LastRevenue, "$"))
				}
			}
			entities.Sim.Mutex.RUnlock()
			if output != "" {
				return output
			}
			built = "Built"
		} else if tile.IsBuildable() {
			built = "Buildable"
		}
		return fmt.Sprintf("Elev: %02d | %s\n%s %s", tile.Elevation, tile.LandSlope, built, tile.LandUse)
	}

	return ""
}
