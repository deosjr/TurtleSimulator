package coords

type Pos struct {
	X, Y, Z int
}

func (p Pos) Add(q Pos) Pos {
	return Pos{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

func (p Pos) Sub(q Pos) Pos {
	return Pos{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

func (p Pos) Up() Pos {
	return Pos{p.X, p.Y, p.Z + 1}
}

func (p Pos) Down() Pos {
	return Pos{p.X, p.Y, p.Z - 1}
}

// some default headings (these are 2D or complex coords)
var (
	North = Pos{0, 1, 0}
	East  = Pos{1, 0, 0}
	South = Pos{0, -1, 0}
	West  = Pos{-1, 0, 0}
)

func HeadingString(p Pos) string {
	switch p {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	}
	return "invalid"
}
