package world

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/janithl/citylyf/internal/entities"
	"github.com/janithl/citylyf/internal/ui/assets"
)

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

	sealevel := entities.Sim.Geography.SeaLevel
	switch tiles[x][y].Elevation {
	case sealevel:
		screen.DrawImage(assets.Assets.Sprites["sand"].Image, op)
	case sealevel - 1:
		if left == sealevel {
			if top == sealevel {
				screen.DrawImage(assets.Assets.Sprites["shore-x-y"].Image, op)
			} else if bottom == sealevel {
				screen.DrawImage(assets.Assets.Sprites["shore-x-y-rev"].Image, op)
			} else if right < sealevel {
				screen.DrawImage(assets.Assets.Sprites["shore-x"].Image, op)
			} else {
				screen.DrawImage(assets.Assets.Sprites["shallowwater"].Image, op)
			}
		} else if right == sealevel {
			if top == sealevel {
				screen.DrawImage(assets.Assets.Sprites["shore-x-rev-y"].Image, op)
			} else if bottom == sealevel {
				screen.DrawImage(assets.Assets.Sprites["shore-x-rev-y-rev"].Image, op)
			} else if left < sealevel {
				screen.DrawImage(assets.Assets.Sprites["shore-x-rev"].Image, op)
			} else {
				screen.DrawImage(assets.Assets.Sprites["shallowwater"].Image, op)
			}
		} else if top == sealevel && bottom < sealevel {
			screen.DrawImage(assets.Assets.Sprites["shore-y"].Image, op)
		} else if bottom == sealevel && top < sealevel {
			screen.DrawImage(assets.Assets.Sprites["shore-y-rev"].Image, op)
		} else {
			screen.DrawImage(assets.Assets.Sprites["shallowwater"].Image, op)
		}
	case 1:
		screen.DrawImage(assets.Assets.Sprites["midwater"].Image, op)
	case 0:
		screen.DrawImage(assets.Assets.Sprites["deepwater"].Image, op)
	default:
		screen.DrawImage(assets.Assets.Sprites["grass"].Image, op)
	}
}

// Renders the mountains
func (wr *WorldRenderer) renderMountains(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
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

	hillLevel := entities.Sim.Geography.HillLevel
	switch tiles[x][y].Elevation {
	case hillLevel + 1:
		screen.DrawImage(assets.Assets.Sprites["mountain"].Image, op)
	case hillLevel:
		screen.DrawImage(assets.Assets.Sprites["hill"].Image, op)
	case hillLevel - 1:
		if left == hillLevel {
			if top == hillLevel {
				screen.DrawImage(assets.Assets.Sprites["slope-x-y"].Image, op)
			} else if bottom == hillLevel {
				screen.DrawImage(assets.Assets.Sprites["slope-x-y-rev"].Image, op)
			} else if right < hillLevel {
				screen.DrawImage(assets.Assets.Sprites["slope-x"].Image, op)
			}
		} else if right == hillLevel {
			if top == hillLevel {
				screen.DrawImage(assets.Assets.Sprites["slope-x-rev-y"].Image, op)
			} else if bottom == hillLevel {
				screen.DrawImage(assets.Assets.Sprites["slope-x-rev-y-rev"].Image, op)
			} else if left < hillLevel {
				screen.DrawImage(assets.Assets.Sprites["slope-x-rev"].Image, op)
			}
		} else if top == hillLevel && bottom < hillLevel {
			screen.DrawImage(assets.Assets.Sprites["slope-y"].Image, op)
		} else if top < hillLevel && bottom == hillLevel {
			screen.DrawImage(assets.Assets.Sprites["slope-y-rev"].Image, op)
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

// Renders retail
func (wr *WorldRenderer) renderRetail(screen *ebiten.Image, op *ebiten.DrawImageOptions, tiles [][]entities.Tile, x, y int) {
	if tiles[x][y].LandUse != entities.RetailUse || tiles[x][y].LandStatus != entities.DevelopedStatus { // not a built shop
		return
	}

	entities.Sim.Mutex.RLock()
	company := entities.Sim.Companies.GetLocationCompany(x, y)
	entities.Sim.Mutex.RUnlock()

	if retailSprite, exists := assets.Assets.Sprites[strings.ToLower(string(company.Industry))+"-small-"+string(company.RoadDirection)]; exists {
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
