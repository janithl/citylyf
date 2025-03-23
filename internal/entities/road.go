package entities

import (
	"fmt"

	"github.com/janithl/citylyf/internal/utils"
)

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

// GetLength returns the road length in number of tiles
func (r *Road) GetLength() int {
	length := 0
	for _, s := range r.Segments {
		length += s.Start.GetDistance(&s.End)
	}
	return length
}

func (r *Road) IsTraversable(x, y int) bool {
	return Sim.Geography.tiles[x][y].Elevation < Sim.Geography.SeaLevel+3
}

func (r *Road) AddSegments(segments []Segment, start bool) {
	if start {
		r.Segments = append(r.Segments, segments...)
	} else {
		r.Segments = append(segments, r.Segments...)
	}
}

func PlaceRoad(start, end Point, roadType RoadType) {
	if start.X == end.Y && start.Y == end.Y {
		return
	}

	var road *Road
	var roadStart bool
	if r, index := Sim.Geography.GetRoadByStartEnd(roadType, start.X, start.Y); r != nil {
		road = r
		roadStart = index == 0
	} else if r, index := Sim.Geography.GetRoadByStartEnd(roadType, end.X, end.Y); r != nil {
		road = r
		roadStart = index == 0
	}

	segments := []Segment{}
	turningPointX, turningPointY := utils.GetTurningPoint(start.X, start.Y, end.X, end.Y)

	if start.X == turningPointX && start.Y != turningPointY {
		segments = append(segments, Segment{
			Start:     Point{start.X, min(start.Y, turningPointY)},
			End:       Point{start.X, max(start.Y, turningPointY)},
			Direction: DirY,
		})
	} else if start.X != turningPointX && start.Y == turningPointY {
		segments = append(segments, Segment{
			Start:     Point{min(start.X, turningPointX), start.Y},
			End:       Point{max(start.X, turningPointX), start.Y},
			Direction: DirX,
		})
	}

	if end.X == turningPointX && end.Y != turningPointY {
		segments = append(segments, Segment{
			Start:     Point{end.X, min(end.Y, turningPointY)},
			End:       Point{end.X, max(end.Y, turningPointY)},
			Direction: DirY,
		})
	} else if end.X != turningPointX && end.Y == turningPointY {
		segments = append(segments, Segment{
			Start:     Point{min(end.X, turningPointX), end.Y},
			End:       Point{max(end.X, turningPointX), end.Y},
			Direction: DirX,
		})
	}

	roadLength := 0
	if road != nil {
		oldLength := road.GetLength()
		road.AddSegments(segments, roadStart)
		roadLength = road.GetLength() - oldLength

		Sim.Geography.placeRoadSegments(segments)
		fmt.Printf("[ Road ] %s extended!\n", road.Name)
	} else {
		road = &Road{
			Name:     Sim.NameService.GetRoadName(),
			Type:     roadType,
			Segments: segments,
		}

		Sim.Geography.addRoad(road)
		roadLength = road.GetLength()
		fmt.Printf("[ Road ] %s opened!\n", road.Name)
	}

	// track road cost
	Sim.Government.AddCapEx(AsphaltRoadConstruction, roadLength)
}
