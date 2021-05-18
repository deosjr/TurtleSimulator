package turtle

import (
    "fmt"

    "github.com/deosjr/TurtleSimulator/blocks"
    "github.com/deosjr/TurtleSimulator/coords"
)

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
    Dig() bool

	// my own functions
	SetProgram(Program)
	Run()
	IsRunning() bool
	String() string
    SetInventory(blocks.Blocktype)
    GetPos() coords.Pos
    GetHeading() coords.Pos
    Tick()
    Tack()
}

type turtle struct {
    blocks.BaseBlock
	pos     coords.Pos
	//heading coords.Pos
	world   *World
	program Program
	tick    chan bool
	ack     chan bool
	running bool
    // hack for now: dont want to build inventory management yet
    inventory blocks.Blocktype
}

type Program func(Turtle)

func (t *turtle) GetPos() coords.Pos {
    return t.pos
}

func (t *turtle) TurnLeft() {
	<-t.tick
	t.Heading = coords.Pos{t.Heading.Y * -1, t.Heading.X, 0}
	t.ack <- true
}

func (t *turtle) TurnRight() {
	<-t.tick
	t.Heading = coords.Pos{t.Heading.Y * 1, -t.Heading.X, 0}
	t.ack <- true
}

// TODO: three lock/unlocks of mutex, can be optimised
func (t *turtle) move(p coords.Pos) error {
    _, ok := t.world.Read(p)
	if ok {
		return fmt.Errorf("block in position")
	}
    t.world.Delete(t.pos)
	t.pos = p
    t.world.Write(p, t)
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
	ok := t.move(t.back())
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

func (t *turtle) detect(p coords.Pos) bool {
	_, ok := t.world.Read(p)
	return ok
}

func (t *turtle) Dig() bool {
	<-t.tick
	ok := t.dig(t.forward())
	t.ack <- true
	return ok
}

func (t *turtle) dig(p coords.Pos) bool {
	_, ok := t.world.Read(p)
    if !ok {
        return false
    }
    t.world.Delete(p)
    return true
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

func (t *turtle) place(p coords.Pos) bool {
	_, ok := t.world.Read(p)
	if ok {
		return false
	}
    toplace := blocks.GetBlock(t.inventory)
    if t.inventory == blocks.Stairs {
	    heading := coords.Pos{t.Heading.Y * -1, t.Heading.X, 0}
        flipped := false
        if t.up() == p {
            if _, upok := t.world.Read(p.Up()); upok {
                flipped = true
            }
        }
        toplace = blocks.BaseBlock{Type: blocks.Stairs, Heading: heading, Flipped: flipped}
    }
    if t.inventory == blocks.CobbleSlab || t.inventory == blocks.BrickSlab {
        flipped := false
        if t.up() == p {
            if _, upok := t.world.Read(p.Up()); upok {
                flipped = true
            }
        }
        toplace = blocks.BaseBlock{Type: t.inventory, Flipped: flipped}
    }
	t.world.Write(p, toplace)
	return true
}

func (t *turtle) SetInventory(bt blocks.Blocktype) {
    t.inventory = bt
}

func (t *turtle) forward() coords.Pos {
	return t.pos.Add(t.Heading)
}

func (t *turtle) back() coords.Pos {
	return t.pos.Sub(t.Heading)
}

func (t *turtle) up() coords.Pos {
    return t.pos.Up()
}

func (t *turtle) down() coords.Pos {
    return t.pos.Down()
}

func (t *turtle) SetProgram(f Program) {
	t.program = f
}

func (t *turtle) Tick() {
    t.tick <- true
}
func (t *turtle) Tack() {
    <-t.ack
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
	switch t.Heading {
	case coords.North:
		return "^"
	case coords.East:
		return ">"
	case coords.South:
		return "v"
	case coords.West:
		return "<"
	}
	return "ERROR"
}

func NewTurtle(p coords.Pos, w *World, heading coords.Pos) Turtle {
	tick := make(chan bool, 1)
	ack := make(chan bool, 1)
	t := &turtle{
        BaseBlock: blocks.BaseBlock{
		    Heading: heading,
            Type: blocks.Turtle,
        },
		pos:     p,
		world:   w,
		tick:    tick,
		ack:     ack,
	}
    w.Write(p, t)
	return t
}
