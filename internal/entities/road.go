package entities

import "fmt"

type IntersectionType string

const (
	NonIntersection IntersectionType = ""
	Roundabout      IntersectionType = "roundabout"
	Fourway         IntersectionType = "fourway"
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
	Asphalt  RoadType = "asphalt"
	Chipseal RoadType = "chipseal"
	Unsealed RoadType = "unsealed"
	NoRoad   RoadType = ""
)

type Road struct {
	Name     string
	Type     RoadType
	Segments []Segment
}

func (r *Road) IsTraversable(x, y int) bool {
	return Sim.Geography.tiles[x][y].Elevation < Sim.Geography.SeaLevel+3
}

func PlaceRoad(start, end Point, roadType RoadType) {
	if start.X == end.Y && start.Y == end.Y {
		return
	}

	segments := []Segment{}
	if start.X != end.X {
		minX, maxX := min(start.X, end.X), max(start.X, end.X)
		segments = append(segments, Segment{
			Start:     Point{minX, start.Y},
			End:       Point{maxX, start.Y},
			Direction: DirX,
		})
		start.X = end.X
	}

	if start.Y != end.Y {
		minY, maxY := min(start.Y, end.Y), max(start.Y, end.Y)
		segments = append(segments, Segment{
			Start:     Point{start.X, minY},
			End:       Point{start.X, maxY},
			Direction: DirY,
		})
	}

	road := &Road{
		Name:     Sim.NameService.GetRoadName(),
		Type:     roadType,
		Segments: segments,
	}

	Sim.Geography.addRoad(road)
	fmt.Printf("[ Road ] %s opened!\n", road.Name)
}
