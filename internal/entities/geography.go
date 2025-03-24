package entities

import (
	"math/rand/v2"

	"github.com/janithl/citylyf/internal/utils"
)

type Geography struct {
	Size, SeaLevel, HillLevel, MaxElevation             int
	peakProbability, rangeProbability, cliffProbability float64
	tiles                                               [][]Tile
	roads                                               []*Road
	Regions                                             Regions
}

// Generate generates the terrain map
// From: https://janithl.github.io/2019/09/go-terrain-gen-part-4/
func (g *Geography) Generate() {
	elevationMap := utils.GenerateElevationMap(g.SeaLevel, g.MaxElevation, g.Size,
		g.peakProbability, g.rangeProbability, g.cliffProbability)

	for x := range elevationMap {
		for y := range elevationMap[x] {
			g.tiles[x][y].Elevation = elevationMap[x][y]
		}
	}
}

// accessor for tiles
func (g *Geography) GetTiles() [][]Tile {
	return g.tiles
}

// bounds check
func (g *Geography) BoundsCheck(x, y int) bool {
	return x >= 0 && y >= 0 && x < g.Size && y < g.Size
}

// accessor for roads
func (g *Geography) GetRoads() []*Road {
	return g.roads
}

func (g *Geography) CheckRoad(x, y int) bool {
	if !g.BoundsCheck(x, y) {
		return false
	}

	return g.tiles[x][y].LandUse == TransportUse
}

// CheckRoadStartEnd returns the first road that starts or ends at x, y
func (g *Geography) GetRoadByStartEnd(roadType RoadType, x, y int) (*Road, int) {
	for _, road := range g.roads {
		segs := len(road.Segments)
		if road.Type != roadType || segs == 0 {
			continue
		}
		if road.Segments[0].Start.X == x && road.Segments[0].Start.Y == y {
			return road, 0
		}
		if road.Segments[segs-1].End.X == x && road.Segments[segs-1].End.Y == y {
			return road, segs - 1
		}
	}
	return nil, 0
}

// zone for land use
func (g *Geography) PlaceLandUse(start Point, end Point, use LandUse) {
	for x := min(start.X, end.X); x <= max(start.X, end.X); x++ {
		for y := min(start.Y, end.Y); y <= max(start.Y, end.Y); y++ {
			roadDir := Sim.Geography.getAccessRoad(x, y)
			if roadDir != "" && g.tiles[x][y].LandUse == NoUse && g.tiles[x][y].IsBuildable() { // zone placeable!
				g.tiles[x][y].LandUse = use
				g.tiles[x][y].LandStatus = UndevelopedStatus
			}
		}
	}
}

// get access road
func (g *Geography) getAccessRoad(x, y int) Direction {
	if !g.BoundsCheck(x, y) || !g.tiles[x][y].IsBuildable() {
		return ""
	}

	if g.CheckRoad(x, y+1) {
		return DirX
	} else if g.CheckRoad(x+1, y) {
		return DirY
	} else if g.CheckRoad(x, y-1) {
		return DirXBack
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
	g.roads = append(g.roads, r)
	g.placeRoadSegments(r.Segments)
}

func (g *Geography) placeRoadSegments(segments []Segment) {
	for _, segment := range segments {
		if segment.Direction == DirX {
			for i := segment.Start.X; i <= segment.End.X; i++ {
				if g.BoundsCheck(i, segment.Start.Y) && g.tiles[i][segment.Start.Y].IsBuildable() && !g.tiles[i][segment.Start.Y].IsBuilt() {
					g.tiles[i][segment.Start.Y].LandUse = TransportUse
					g.tiles[i][segment.Start.Y].LandStatus = DevelopedStatus
				}
				g.setIntersectionType(i-1, segment.Start.Y)
			}
			g.setIntersectionType(segment.End.X, segment.Start.Y)
		} else if segment.Direction == DirY {
			for i := segment.Start.Y; i <= segment.End.Y; i++ {
				if g.BoundsCheck(segment.Start.X, i) && g.tiles[segment.Start.X][i].IsBuildable() && !g.tiles[segment.Start.X][i].IsBuilt() {
					g.tiles[segment.Start.X][i].LandUse = TransportUse
					g.tiles[segment.Start.X][i].LandStatus = DevelopedStatus
				}
				g.setIntersectionType(segment.Start.X, i-1)
			}
			g.setIntersectionType(segment.Start.X, segment.End.Y)
		}
	}
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

// get potential site to place a building
func (g *Geography) GetPotentialSite(use LandUse) *Point {
	potentialSites := []*Point{}
	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			if g.tiles[x][y].LandUse == use && Sim.Geography.tiles[x][y].LandStatus == UndevelopedStatus {
				roadDir := g.getAccessRoad(x, y)
				if roadDir == "" {
					continue
				}
				potentialSites = append(potentialSites, &Point{X: x, Y: y})
			}
		}
	}

	if len(potentialSites) < 1 {
		return nil
	}

	return potentialSites[rand.IntN(len(potentialSites))]
}

// NewGeography returns a new terrain map
func NewGeography(mapSize, regionSize, maxElevation, SeaLevel, HillLevel int, peakProbability, rangeProbability, cliffProbability float64) *Geography {
	tiles := make([][]Tile, mapSize)
	for i := 0; i < mapSize; i++ {
		tiles[i] = make([]Tile, mapSize)
	}

	geography := &Geography{
		Size:             mapSize,
		MaxElevation:     maxElevation,
		SeaLevel:         SeaLevel,
		HillLevel:        HillLevel,
		peakProbability:  peakProbability,
		rangeProbability: rangeProbability,
		cliffProbability: cliffProbability,
		tiles:            tiles,
		Regions:          NewRegions(mapSize, regionSize),
	}
	// generate the terrain
	geography.Generate()

	return geography
}
