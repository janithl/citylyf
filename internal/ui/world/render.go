package world

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/assets"
)

// Renders the base tile
func (wr *WorldRenderer) renderBaseTiles(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	switch tiles[x][y].Elevation {
	case entities.Sim.Geography.SeaLevel:
		if sprite, exists := assets.Assets.Sprites[string(tiles[x][y].LandSlope)+"-sand"]; exists {
			screen.DrawImage(sprite.Image, op)
		} else {
			screen.DrawImage(assets.Assets.Sprites["flat-sand"].Image, op)
		}
	case entities.Sim.Geography.SeaLevel - 1:
		screen.DrawImage(assets.Assets.Sprites["shallow-water"].Image, op)
	case 1:
		screen.DrawImage(assets.Assets.Sprites["mid-water"].Image, op)
	case 0:
		screen.DrawImage(assets.Assets.Sprites["deep-water"].Image, op)
	default:
		screen.DrawImage(assets.Assets.Sprites["flat-grass"].Image, op)
	}

	// draw tile borders
	borderOp := *op
	borderOp.ColorScale.Scale(1, 1, 1, 0.4)
	if tiles[x][y].Elevation < entities.Sim.Geography.SeaLevel {
		screen.DrawImage(assets.Assets.Sprites["ui-tile-border"].Image, &borderOp)
	}
}

// Renders the mountains
func (wr *WorldRenderer) renderMountains(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	hillLevel := entities.Sim.Geography.HillLevel
	switch tiles[x][y].Elevation {
	case hillLevel + 1:
		mountainOp := *op
		mountainOp.GeoM.Translate(0, wr.elevationToZ(tiles[x][y].Elevation)*wr.zoomFactor)
		screen.DrawImage(assets.Assets.Sprites["mountain-peak"].Image, &mountainOp)
	case hillLevel:
		if sprite, exists := assets.Assets.Sprites[string(tiles[x][y].LandSlope)+"-hill"]; exists {
			screen.DrawImage(sprite.Image, op)
		} else {
			screen.DrawImage(assets.Assets.Sprites["flat-hill"].Image, op)
		}
	}
}

// Renders houses
func (wr *WorldRenderer) renderHouses(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	if tiles[x][y].LandUse != entities.ResidentialUse || tiles[x][y].LandStatus != entities.DevelopedStatus { // not a built house
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

// Renders industries/shops
func (wr *WorldRenderer) renderIndusty(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	if !(tiles[x][y].LandUse == entities.RetailUse || tiles[x][y].LandUse == entities.AgricultureUse) || tiles[x][y].LandStatus != entities.DevelopedStatus { // not a built shop or farm
		return
	}

	entities.Sim.Mutex.RLock()
	company := entities.Sim.Companies.GetLocationCompany(x, y)
	entities.Sim.Mutex.RUnlock()

	if company == nil {
		return
	}

	if retailSprite, exists := assets.Assets.Sprites["industry-"+strings.ToLower(string(company.Industry))+"-"+string(company.RoadDirection)]; exists {
		screen.DrawImage(retailSprite.Image, op)
	}
}

// Renders roads
func (wr *WorldRenderer) renderRoads(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	if tiles[x][y].LandUse != entities.TransportUse { // not a road
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

func (wr *WorldRenderer) assignAnimations() {
	for _, region := range entities.Sim.Geography.Regions {
		delay := 0
		for _, trip := range region.Trips {
			if trip.Start == nil || trip.End == nil {
				continue
			}

			for _, anim := range wr.animations {
				if anim.IsFinished() {
					anim.SetPath(entities.Sim.Geography.FindPath(trip.Start, trip.End))
					anim.CalculateSpeed(delay)
					delay += 60 // delay next animation by 1 seconds
					break
				}
			}
		}
	}
}
