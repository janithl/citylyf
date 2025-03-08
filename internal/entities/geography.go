package entities

import (
	"math/rand"
)

type Tile struct {
	Elevation    int
	Road, House  bool
	Intersection IntersectionType
	Zone         Zone
}

type Geography struct {
	Size, SeaLevel, MaxElevation                        int
	biasX, biasY                                        int // bias x and y create a vector along which mountain ranges form
	peakProbability, rangeProbability, cliffProbability float64
	tiles                                               [][]Tile
	roads                                               []Road
}

// Generate generates the terrain map
// From: https://janithl.github.io/2019/09/go-terrain-gen-part-4/
func (g *Geography) Generate() {
	// iterate down from max elevation, assigning vals
	for e := g.MaxElevation; e > 0; e-- {
		for y := 0; y < g.Size; y++ {
			for x := 0; x < g.Size; x++ {
				// if the element is next to a element with elevation x, it
				// should get elevation x - 1
				// alternately, if the random value meets our criteria, it's a peak
				if g.adjacentElevation(x, y, e) || rand.Float64() < g.peakProbability {
					g.setElevation(x, y, e)
					if rand.Float64() > g.rangeProbability { // randomly add follow-up peaks
						g.setElevation(x+g.biasX, y+g.biasY, e)
					}
					if rand.Float64() > g.rangeProbability {
						g.setElevation(x-g.biasX, y-g.biasY, e)
					}
				}
			}
		}
	}
}

// adjacentElevation checks if an adjacent element
// to the given element (h, w) is at a given elevation
func (g *Geography) adjacentElevation(w, h, elevation int) bool {
	adjacentTiles := []Point{
		{w - 1, h - 1}, {w, h - 1}, {w + 1, h - 1},
		{w - 1, h}, {w + 1, h},
		{w - 1, h + 1}, {w, h + 1}, {w + 1, h + 1}}

	for _, tile := range adjacentTiles {
		if g.BoundsCheck(tile.X, tile.Y) && g.tiles[tile.X][tile.Y].Elevation == elevation+1 {
			// if this element is *not* randomly a cliff, return true
			return rand.Float64() > g.cliffProbability
		}
	}

	return false
}

// setElevation checks if a tile exists and updates the elevation
// to the given element (h, w) is at a given elevation
func (g *Geography) setElevation(x, y, elevation int) {
	if !g.BoundsCheck(x, y) {
		return
	}

	// if the element has already been assigned, skip it
	if g.tiles[x][y].Elevation > 0 {
		return
	}

	g.tiles[x][y].Elevation = elevation
}

// accessor for tiles
func (g *Geography) GetTiles() [][]Tile {
	return g.tiles
}

// accessor for roads
func (g *Geography) GetRoads() []Road {
	return g.roads
}

func (g *Geography) CheckRoad(x, y int) bool {
	if !g.BoundsCheck(x, y) {
		return false
	}

	return g.tiles[x][y].Road
}

// place residential zone
func (g *Geography) PlaceZone(start Point, end Point, zone Zone) {
	for x := min(start.X, end.X); x <= max(start.X, end.X); x++ {
		for y := min(start.Y, end.Y); y <= max(start.Y, end.Y); y++ {
			if roadDir := Sim.Geography.getAccessRoad(x, y); roadDir != "" { // zone placeable!
				g.tiles[x][y].Zone = zone
			}
		}
	}
}

// get access road
func (g *Geography) getAccessRoad(x, y int) Direction {
	if !g.BoundsCheck(x, y) || g.tiles[x][y].Elevation < g.SeaLevel || g.tiles[x][y].Road || g.tiles[x][y].House {
		return ""
	}

	if g.CheckRoad(x, y+1) {
		return DirX
	} else if g.CheckRoad(x, y-1) {
		return DirXBack
	} else if g.CheckRoad(x+1, y) {
		return DirY
	} else if g.CheckRoad(x-1, y) {
		return DirYBack
	}

	return ""
}

// setIntersectionType sets the type of intersection based on surrounding tiles
func (g *Geography) setIntersectionType(x, y int) {
	if !g.CheckRoad(x, y) {
		return
	}

	top := g.CheckRoad(x, y+1)
	bottom := g.CheckRoad(x, y-1)
	right := g.CheckRoad(x+1, y)
	left := g.CheckRoad(x-1, y)
	switch {
	case top && bottom && left && right:
		g.tiles[x][y].Intersection = Fourway
	case top && bottom && left:
		g.tiles[x][y].Intersection = ThreewayYUp
	case top && bottom && right:
		g.tiles[x][y].Intersection = ThreewayYDown
	case top && left && right:
		g.tiles[x][y].Intersection = ThreewayXDown
	case bottom && left && right:
		g.tiles[x][y].Intersection = ThreewayXUp
	case top && left:
		g.tiles[x][y].Intersection = BendTopLeft
	case top && right:
		g.tiles[x][y].Intersection = BendTopRight
	case bottom && left:
		g.tiles[x][y].Intersection = BendBottomLeft
	case bottom && right:
		g.tiles[x][y].Intersection = BendBottomRight
	default:
		g.tiles[x][y].Intersection = NonIntersection
	}
}

// add a new road
func (g *Geography) addRoad(r *Road) {
	g.roads = append(g.roads, *r)
	for _, segment := range r.Segments {
		if segment.Direction == DirX {
			for i := segment.Start.X; i <= segment.End.X; i++ {
				if g.BoundsCheck(i, segment.Start.Y) && !g.tiles[i][segment.Start.Y].House {
					g.tiles[i][segment.Start.Y].Road = true
					g.tiles[i][segment.Start.Y].Zone = ""
				}
				g.setIntersectionType(i-1, segment.Start.Y)
			}
			g.setIntersectionType(segment.End.X, segment.Start.Y)
		} else if segment.Direction == DirY {
			for i := segment.Start.Y; i <= segment.End.Y; i++ {
				if g.BoundsCheck(segment.Start.X, i) && !g.tiles[segment.Start.X][i].House {
					g.tiles[segment.Start.X][i].Road = true
					g.tiles[segment.Start.X][i].Zone = ""
				}
				g.setIntersectionType(segment.Start.X, i-1)
			}
			g.setIntersectionType(segment.Start.X, segment.End.Y)
		}
	}
}

// bounds check
func (g *Geography) BoundsCheck(x, y int) bool {
	return x >= 0 && y >= 0 && x < g.Size && y < g.Size
}

// check if coordinates are within a road segment, and that road's direction and type
func (g *Geography) IsWithinRoad(x, y int) (Direction, RoadType) {
	for _, road := range g.roads {
		for _, segment := range road.Segments {
			if segment.Direction == DirX && y == segment.Start.Y && x >= segment.Start.X && x <= segment.End.X {
				return DirX, road.Type
			}
			if segment.Direction == DirY && x == segment.Start.X && y >= segment.Start.Y && y <= segment.End.Y {
				return DirY, road.Type
			}
		}
	}
	return "", ""
}

// toggle roundabout
func (g *Geography) ToggleRoundabout(x, y int) {
	if !g.BoundsCheck(x, y) {
		return
	}

	if g.tiles[x][y].Intersection == Fourway {
		g.tiles[x][y].Intersection = Roundabout
	} else if g.tiles[x][y].Intersection == Roundabout {
		g.tiles[x][y].Intersection = Fourway
	}
}

// NewGeography returns a new terrain map
func NewGeography(size, maxElevation, SeaLevel int, peakProbability, rangeProbability, cliffProbability float64) *Geography {
	tiles := make([][]Tile, size)
	for i := 0; i < size; i++ {
		tiles[i] = make([]Tile, size)
	}

	geography := &Geography{
		Size:             size,
		MaxElevation:     maxElevation,
		SeaLevel:         SeaLevel,
		biasX:            rand.Intn(6) - 3,
		biasY:            rand.Intn(6) - 3,
		peakProbability:  peakProbability,
		rangeProbability: rangeProbability,
		cliffProbability: cliffProbability,
		tiles:            tiles,
	}
	// generate the terrain
	geography.Generate()

	return geography
}
