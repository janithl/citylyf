package entities

type Point struct {
	X, Y int
}

func (p *Point) Equal(q *Point) bool {
	return p.X == q.X && p.Y == q.Y
}
