package geo

import "github.com/janithl/citylyf/internal/entities"

// FindPath uses BFS to find a path from source to dest using only tiles with Road == true.
func FindPath(source, dest entities.Point) []entities.Point {
	tiles := entities.Sim.Geography.GetTiles()

	// Check that source and destination are within bounds.
	if !entities.Sim.Geography.BoundsCheck(source.X, source.Y) || !entities.Sim.Geography.BoundsCheck(dest.X, dest.Y) {
		return nil
	}

	// Ensure source and destination are road tiles.
	if !tiles[source.X][source.Y].Road || !tiles[dest.X][dest.Y].Road {
		return nil
	}

	// Directions: up, down, left, right.
	directions := []entities.Point{
		{X: 0, Y: 1},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: -1, Y: 0},
	}

	// Maps to keep track of visited tiles and how we got there.
	visited := make(map[entities.Point]bool)
	cameFrom := make(map[entities.Point]entities.Point)

	queue := []entities.Point{source}
	visited[source] = true

	found := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == dest {
			found = true
			break
		}

		for _, d := range directions {
			neighbor := entities.Point{X: current.X + d.X, Y: current.Y + d.Y}
			if !entities.Sim.Geography.BoundsCheck(neighbor.X, neighbor.Y) {
				continue
			}
			if visited[neighbor] {
				continue
			}
			// Only consider road tiles.
			if !tiles[neighbor.X][neighbor.Y].Road {
				continue
			}
			visited[neighbor] = true
			cameFrom[neighbor] = current
			queue = append(queue, neighbor)
		}
	}

	if !found {
		return nil
	}

	// Reconstruct path by walking backward from dest to source.
	var path []entities.Point
	for cur := dest; cur != source; {
		path = append([]entities.Point{cur}, path...)
		cur = cameFrom[cur]
	}
	path = append([]entities.Point{source}, path...)

	return path
}
