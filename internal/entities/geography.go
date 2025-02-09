package entities

import (
	"math/rand"
)

type Tile struct {
	Elevation    int
	Road         bool
	Intersection bool
	Roundabout   bool
}

type Geography struct {
	Size, SeaLevel, MaxElevation      int
	peakProbability, cliffProbability float64
	tiles                             [][]Tile
	roads                             []Road
}

// Generate generates the terrain map
// From: https://janithl.github.io/2019/09/go-terrain-gen-part-4/
func (g *Geography) Generate() {
	// iterate down from max elevation, assigning vals
	for e := g.MaxElevation; e > 0; e-- {
		for h := 0; h < g.Size; h++ {
			for w := 0; w < g.Size; w++ {
				// if the element has already been assigned, skip it
				if g.tiles[h][w].Elevation > 0 {
					continue
				}

				// if the element is next to a element with elevation x, it
				// should get elevation x - 1
				// alternately, if the random value meets our criteria, it's a peak
				if g.adjacentElevation(h, w, e) || rand.Float64() < g.peakProbability {
					g.tiles[h][w].Elevation = e
				}
			}
		}
	}
}

// adjacentElevation checks if an adjacent element
// to the given element (h, w) is at a given elevation
func (g *Geography) adjacentElevation(h, w, elevation int) bool {
	for y := max(0, h-1); y <= min(g.Size-1, h+1); y++ {
		for x := max(0, w-1); x <= min(g.Size-1, w+1); x++ {
			if g.tiles[y][x].Elevation == elevation+1 {
				// if this element is *not* randomly a cliff, return true
				return rand.Float64() > g.cliffProbability
			}
		}
	}

	return false
}

// accessor for tiles
func (g *Geography) GetTiles() [][]Tile {
	return g.tiles
}

// accessor for roads
func (g *Geography) GetRoads() []Road {
	return g.roads
}

// add a new road
func (g *Geography) addRoad(r *Road) {
	g.roads = append(g.roads, *r)
	for _, segment := range r.Segments {
		if segment.Direction == DirX {
			for i := segment.Start.X; i <= segment.End.X; i++ {
				if g.tiles[i][segment.Start.Y].Road {
					g.tiles[i][segment.Start.Y].Intersection = true
				} else {
					g.tiles[i][segment.Start.Y].Road = true
				}
			}
		} else if segment.Direction == DirY {
			for i := segment.Start.Y; i <= segment.End.Y; i++ {
				if g.tiles[segment.Start.X][i].Road {
					g.tiles[segment.Start.X][i].Intersection = true
				} else {
					g.tiles[segment.Start.X][i].Road = true
				}
			}
		}
	}
}

// bounds check
func (g *Geography) BoundsCheck(x, y int) bool {
	return x >= 0 && y >= 0 && x < Sim.Geography.Size && y < Sim.Geography.Size
}

// check if coordinates are within a road segment, and that road's direction
func (g *Geography) IsWithinRoad(x, y int) Direction {
	for _, road := range g.roads {
		for _, segment := range road.Segments {
			if segment.Direction == DirX && y == segment.Start.Y && x >= segment.Start.X && x <= segment.End.X {
				return DirX
			}
			if segment.Direction == DirY && x == segment.Start.X && y >= segment.Start.Y && y <= segment.End.Y {
				return DirY
			}
		}
	}
	return ""
}

// NewGeography returns a new terrain map
func NewGeography(size, maxElevation, SeaLevel int, peakProbability, cliffProbability float64) *Geography {
	tiles := make([][]Tile, size)
	for i := 0; i < size; i++ {
		tiles[i] = make([]Tile, size)
	}

	geography := &Geography{
		Size:             size,
		MaxElevation:     maxElevation,
		SeaLevel:         SeaLevel,
		peakProbability:  peakProbability,
		cliffProbability: cliffProbability,
		tiles:            tiles,
	}
	geography.Generate()

	return geography
}
