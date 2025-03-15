package entities

type Point struct {
	X, Y int
}

func (p *Point) Equal(q *Point) bool {
	return p.X == q.X && p.Y == q.Y
}

func (p *Point) GetNeighbours(distance int, cardinal bool) []*Point {
	if cardinal {
		return []*Point{
			{X: p.X - distance, Y: p.Y},
			{X: p.X, Y: p.Y - distance},
			{X: p.X, Y: p.Y + distance},
			{X: p.X + distance, Y: p.Y},
		}
	}
	return []*Point{
		{X: p.X - distance, Y: p.Y - distance},
		{X: p.X - distance, Y: p.Y},
		{X: p.X - distance, Y: p.Y + distance},
		{X: p.X, Y: p.Y - distance},
		{X: p.X, Y: p.Y + distance},
		{X: p.X + distance, Y: p.Y - distance},
		{X: p.X + distance, Y: p.Y},
		{X: p.X + distance, Y: p.Y + distance},
	}
}
