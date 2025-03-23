package entities

import (
	"testing"
)

func TestRoadLength(t *testing.T) {
	// create a road with segments with a length of 10
	road := Road{
		Segments: []Segment{
			{
				Start:     Point{0, 10},
				End:       Point{0, 5},
				Direction: DirY,
			},
			{
				Start:     Point{0, 5},
				End:       Point{5, 5},
				Direction: DirX,
			},
		},
	}

	expected := 10
	if length := road.GetLength(); length != expected {
		t.Errorf("Expected %d for road length, got %d", expected, length)
	}

	// add segments with a length of 10 more
	road.AddSegments([]Segment{
		{
			Start:     Point{5, 5},
			End:       Point{10, 5},
			Direction: DirX,
		},
		{
			Start:     Point{10, 5},
			End:       Point{10, 10},
			Direction: DirY,
		},
	}, false)

	expected = 20
	if length := road.GetLength(); length != expected {
		t.Errorf("Expected %d for road length, got %d", expected, length)
	}
}
