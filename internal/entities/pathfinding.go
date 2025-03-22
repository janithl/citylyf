package entities

// FindPath uses BFS to find a path from source to dest using only tiles with Road == true.
func (g *Geography) FindPath(source, dest *Point) []*Point {
	// Check that source and destination are within bounds.
	if !g.BoundsCheck(source.X, source.Y) || !g.BoundsCheck(dest.X, dest.Y) {
		return nil
	}

	// Ensure source and destination are road tiles.
	if g.tiles[source.X][source.Y].LandUse != TransportUse || g.tiles[dest.X][dest.Y].LandUse != TransportUse {
		return nil
	}

	// Maps to keep track of visited tiles and how we got there.
	visited := make(map[Point]bool)
	cameFrom := make(map[Point]*Point)

	queue := []*Point{source}
	visited[*source] = true

	found := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Equal(dest) {
			found = true
			break
		}

		for _, neighbor := range current.GetNeighbours(1, true) {
			if !g.BoundsCheck(neighbor.X, neighbor.Y) {
				continue
			}
			if visited[*neighbor] {
				continue
			}
			// Only consider road tiles.
			if g.tiles[neighbor.X][neighbor.Y].LandUse != TransportUse {
				continue
			}
			visited[*neighbor] = true
			cameFrom[*neighbor] = current
			queue = append(queue, neighbor)
		}
	}

	if !found {
		return nil
	}

	// Reconstruct path by walking backward from dest to source.
	var path []*Point
	for cur := dest; cur != source; {
		path = append([]*Point{cur}, path...)
		cur = cameFrom[*cur]
	}
	path = append([]*Point{source}, path...)

	return path
}

// Finds the turning points along the path
func (g *Geography) FindTurns(path []*Point) []*Point {
	if len(path) < 3 {
		return nil // No turns in a straight line with <3 points
	}

	var turns []*Point

	// Get initial direction
	prevDir := &Point{X: path[1].X - path[0].X, Y: path[1].Y - path[0].Y}

	for i := 1; i < len(path)-1; i++ {
		// Compute current direction
		curDir := &Point{X: path[i+1].X - path[i].X, Y: path[i+1].Y - path[i].Y}

		// If direction changes, it's a turn
		if !curDir.Equal(prevDir) {
			turns = append(turns, path[i])
		}

		prevDir = curDir
	}

	return turns
}
