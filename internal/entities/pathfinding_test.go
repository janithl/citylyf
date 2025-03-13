package entities

import "testing"

func TestFindPath(t *testing.T) {
	// Create a mock tile grid
	size := 8
	tiles := make([][]Tile, size)
	for i := range tiles {
		tiles[i] = make([]Tile, size)
	}

	// Define a road layout (simple straight road from (0,0) to (4,0))
	for x := 0; x < size; x++ {
		tiles[x][0].Road = true
	}

	g := &Geography{
		tiles: tiles,
		Size:  size,
	}

	t.Run("FindPath", func(t *testing.T) {
		source := &Point{X: 0, Y: 0}
		dest := &Point{X: 7, Y: 0}
		path := g.FindPath(source, dest)
		if path == nil {
			t.Errorf("Expected a valid path, got nil")
			return
		}
		if len(path) != 8 {
			t.Errorf("Expected path length of 8, got %d", len(path))
		}
		if path[0] != source {
			t.Errorf("Expected first point in path to be %v, got %v", source, path[0])
		}
		if path[len(path)-1] != dest {
			t.Errorf("Expected last point in path to be %v, got %v", dest, path[len(path)-1])
		}
	})

	// Test for a case with no valid path (road disconnected)
	t.Run("FindPath: Non road", func(t *testing.T) {
		nonRoadSource := &Point{X: 0, Y: 1} // Not a road
		dest := &Point{X: 7, Y: 0}
		path := g.FindPath(nonRoadSource, dest)
		if path != nil {
			t.Errorf("Expected nil path for non-road start, got %v", path)
		}
	})
}

func TestFindTurns(t *testing.T) {
	g := &Geography{}

	// Helper function to create a Point slice
	newPath := func(points ...[2]int) []*Point {
		var path []*Point
		for _, p := range points {
			path = append(path, &Point{X: p[0], Y: p[1]})
		}
		return path
	}

	// Test case: No turns (straight horizontal)
	path1 := newPath([2]int{0, 0}, [2]int{1, 0}, [2]int{2, 0}, [2]int{3, 0})
	if turns := g.FindTurns(path1); len(turns) != 0 {
		t.Errorf("Expected no turns, got %d", len(turns))
	}

	// Test case: One turn (right angle)
	path2 := newPath([2]int{0, 0}, [2]int{1, 0}, [2]int{1, 1}, [2]int{1, 2})
	expectedTurns := 1
	if turns := g.FindTurns(path2); len(turns) != expectedTurns {
		t.Errorf("Expected %d turn, got %d", expectedTurns, len(turns))
	}

	// Test case: Two turns (zigzag)
	path3 := newPath([2]int{0, 0}, [2]int{1, 0}, [2]int{1, 1}, [2]int{2, 1}, [2]int{3, 1})
	expectedTurns = 2
	if turns := g.FindTurns(path3); len(turns) != expectedTurns {
		t.Errorf("Expected %d turns, got %d", expectedTurns, len(turns))
	}

	// Test case: Too short to turn
	path4 := newPath([2]int{0, 0}, [2]int{1, 0})
	if turns := g.FindTurns(path4); turns != nil {
		t.Errorf("Expected nil for short path, got %v", turns)
	}
}
