package entities

import (
	"fmt"
	"math/rand"
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

func NewRoad(startX, startY, endX, endY int, roadType RoadType) *Road {
	road := &Road{
		Name: fmt.Sprintf("Street %d", rand.Intn(100)),
		Type: roadType,
		Segments: []Segment{
			{
				Start:     Point{startX, startY},
				End:       Point{endX, startY},
				Direction: DirX,
			},
			{
				Start:     Point{endX, startY},
				End:       Point{endX, endY},
				Direction: DirY,
			},
		},
	}

	return road
}
