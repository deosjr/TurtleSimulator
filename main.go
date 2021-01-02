package main

type pos struct {
	x, y, z int
}

func (p pos) add(q pos) pos {
	return pos{p.x + q.x, p.y + q.y, p.z + q.z}
}

func (p pos) sub(q pos) pos {
	return pos{p.x - q.x, p.y - q.y, p.z - q.z}
}

// convention: world minimal coord is 0,0,0 and max dim,dim,dim
type world struct {
	grid map[pos]Block
	dim  int
    turtles []Turtle
    tick chan bool
    tack chan bool
}

func NewWorld(dim int, tick, tack chan bool) *world {
    return &world{
        grid: map[pos]Block{},
        dim: dim,
        tick: tick,
        tack: tack,
    }
}

type Block interface{}
type block struct{}
type grass struct{}
type stone struct{}

func wallbuildfunc() program {
    return func(t Turtle) {
		sidestepRight := func() {
			t.TurnRight()
			t.Forward()
			t.TurnLeft()
		}
		sidestepLeft := func() {
			t.TurnLeft()
			t.Forward()
			t.TurnRight()
		}
		if t.Detect() {
			sidestepRight()
			t.Forward()
			t.TurnLeft()
		} else {
			t.Forward()
			t.TurnLeft()
			t.Back()
		}
		for i := 0; i < 2; i++ {
			t.Place()
			t.TurnRight()
			t.Place()
			for j := 0; j < 2; j++ {
				sidestepRight()
				t.Place()
			}
			t.Up()
			t.Place()
			for j := 0; j < 2; j++ {
				sidestepLeft()
				t.Place()
			}
			t.TurnLeft()
			t.Place()
			t.Up()
		}
		t.Place()
		t.Up()
		t.Place()
		t.Down()
		t.TurnRight()
		t.Place()
		t.PlaceUp()

		for i := 0; i < 2; i++ {
			sidestepRight()
			t.Place()
			t.PlaceUp()
		}
		t.Down()
		t.TurnRight()
		t.PlaceUp()
		t.Back()
		t.Back()
		t.TurnRight()
		t.TurnRight()
		t.PlaceUp()
		for i := 0; i < 3; i++ {
			t.Down()
			t.Back()
		}
		t.TurnRight()
		t.Back()
    }
}

// TODO: a game tick in minecraft is 1/20 second
// from dan200 computercraft java code:
// each animation takes 8 ticks to complete unless otherwise specified.

func main() {
	tick := make(chan bool, 1)
	tack := make(chan bool, 1)
    w := NewWorld(5, tick, tack)
	//w.grid[pos{0, 3, 0}] = block{}
	//w.grid[pos{2, 2, 0}] = block{}
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			w.grid[pos{x, y, -1}] = grass{}
		}
	}
	t := NewTurtle(pos{0, 0, 0}, w, tick, tack)
    w.turtles = []Turtle{t}
	// wall building program
	t.SetProgram(wallbuildfunc())

    //vis := NewRaytracer()
    vis := ascii{}
    Visualise(vis, w)
}
