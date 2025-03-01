package entities

import (
	"fmt"
	"math/rand"
)

type IntersectionType string

const (
	NonIntersection IntersectionType = ""
	Intersection    IntersectionType = "intersection"
	ThreewayXUp     IntersectionType = "threeway-x-up"
	ThreewayXDown   IntersectionType = "threeway-x-down"
	ThreewayYUp     IntersectionType = "threeway-y-up"
	ThreewayYDown   IntersectionType = "threeway-y-down"
	BendTopLeft     IntersectionType = "bend-top-left"
	BendTopRight    IntersectionType = "bend-top-right"
	BendBottomLeft  IntersectionType = "bend-bottom-left"
	BendBottomRight IntersectionType = "bend-bottom-right"
)

type Segment struct {
	Start, End Point
	Direction  Direction
}

type RoadType string

const (
	Asphalt  RoadType = "Asphalt"
	Chipseal RoadType = "Chipseal"
	Unsealed RoadType = "Unsealed"
)

type Road struct {
	Name     string
	Type     RoadType
	Segments []Segment
}

func (r *Road) IsTraversable(x, y int) bool {
	return Sim.Geography.tiles[x][y].Elevation < Sim.Geography.SeaLevel+3
}

func PlaceRoad(startX, startY, endX, endY int, roadType RoadType) {
	if startX == endY && startY == endY {
		return
	}

	segments := []Segment{}
	if startX != endX {
		minX, maxX := min(startX, endX), max(startX, endX)
		segments = append(segments, Segment{
			Start:     Point{minX, startY},
			End:       Point{maxX, startY},
			Direction: DirX,
		})
		startX = endX
	}

	if startY != endY {
		minY, maxY := min(startY, endY), max(startY, endY)
		segments = append(segments, Segment{
			Start:     Point{startX, minY},
			End:       Point{startX, maxY},
			Direction: DirY,
		})
	}

	road := &Road{
		Name:     fmt.Sprintf("Street %d", rand.Intn(100)),
		Type:     roadType,
		Segments: segments,
	}

	Sim.Geography.addRoad(road)
}
