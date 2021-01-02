package main

import "fmt"

// use interface to mimic how we can't access everything in CC:Tweaked API
type Turtle interface {
	// CC:Tweaked API
	TurnLeft()
	TurnRight()
	Forward() error
	Back() error
	Up() error
	Down() error
	Detect() bool
	DetectUp() bool
	DetectDown() bool
	Place() bool
	PlaceUp() bool
	PlaceDown() bool

	// my own functions
	SetProgram(program)
	Run()
	IsRunning() bool
	String() string
}

type turtle struct {
	pos     pos
	heading pos
	world   *world
	program program
	tick    <-chan bool
	ack     chan<- bool
	running bool
}

type program func(Turtle)

func (t *turtle) TurnLeft() {
	<-t.tick
	t.heading = pos{t.heading.y * -1, t.heading.x, 0}
	t.ack <- true
}

func (t *turtle) TurnRight() {
	<-t.tick
	t.heading = pos{t.heading.y * 1, -t.heading.x, 0}
	t.ack <- true
}

func (t *turtle) move(p pos) error {
	_, ok := t.world.grid[p]
	if ok {
		return fmt.Errorf("block in position")
	}
	delete(t.world.grid, t.pos)
	t.pos = p
	t.world.grid[p] = t
	return nil
}

func (t *turtle) Forward() error {
	<-t.tick
	ok := t.move(t.forward())
	t.ack <- true
	return ok
}

func (t *turtle) Back() error {
	<-t.tick
	newpos := t.pos.sub(t.heading)
	ok := t.move(newpos)
	t.ack <- true
	return ok
}

func (t *turtle) Up() error {
	<-t.tick
	ok := t.move(t.up())
	t.ack <- true
	return ok
}

func (t *turtle) Down() error {
	<-t.tick
	ok := t.move(t.down())
	t.ack <- true
	return ok
}

func (t *turtle) Detect() bool {
	<-t.tick
	ok := t.detect(t.forward())
	t.ack <- true
	return ok
}

func (t *turtle) DetectUp() bool {
	<-t.tick
	ok := t.detect(t.up())
	t.ack <- true
	return ok
}

func (t *turtle) DetectDown() bool {
	<-t.tick
	ok := t.detect(t.down())
	t.ack <- true
	return ok
}

func (t *turtle) detect(p pos) bool {
	_, ok := t.world.grid[p]
	return ok
}

func (t *turtle) Place() bool {
	<-t.tick
	ok := t.place(t.forward())
	t.ack <- true
	return ok
}

func (t *turtle) PlaceUp() bool {
	<-t.tick
	ok := t.place(t.up())
	t.ack <- true
	return ok
}

func (t *turtle) PlaceDown() bool {
	<-t.tick
	ok := t.place(t.down())
	t.ack <- true
	return ok
}

func (t *turtle) place(p pos) bool {
	_, ok := t.world.grid[p]
	if ok {
		return false
	}
	t.world.grid[p] = block{}
	return true
}

func (t *turtle) forward() pos {
	return t.pos.add(t.heading)
}

func (t *turtle) up() pos {
	return pos{t.pos.x, t.pos.y, t.pos.z + 1}
}

func (t *turtle) down() pos {
	return pos{t.pos.x, t.pos.y, t.pos.z - 1}
}

func (t *turtle) SetProgram(f program) {
	t.program = f
}

func (t *turtle) Run() {
	t.running = true
	t.program(t)
	t.running = false
}

func (t *turtle) IsRunning() bool {
	return t.running
}

func (t *turtle) String() string {
	switch t.heading {
	case pos{0, 1, 0}:
		return "^"
	case pos{-1, 0, 0}:
		return "<"
	case pos{0, -1, 0}:
		return "v"
	case pos{1, 0, 0}:
		return ">"
	}
	return "ERROR"
}

// turtle starts with heading north
func NewTurtle(p pos, w *world, tick <-chan bool, ack chan<- bool) Turtle {
	t := &turtle{
		pos:     p,
		heading: pos{0, 1, 0},
		world:   w,
		tick:    tick,
		ack:     ack,
	}
	w.grid[p] = t
	return t
}
