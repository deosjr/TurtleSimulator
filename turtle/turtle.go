package turtle

import (
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
	Inspect() (blocks.Block, bool)

	// my own functions
	SetProgram(Program)
	Run()
	IsRunning() bool
	String() string
	GetPos() coords.Pos
	GetHeading() coords.Pos
	Tick()
	Tack()
	// DEBUG
	SetInventory(blocks.Blocktype)
	SetInfiniteInventory()
	SetInfiniteFuel()
}

type turtle struct {
	blocks.BaseBlock
	pos     coords.Pos
	world   *World
	program Program
	tick    chan bool
	ack     chan bool
	running bool
	// inventory management
	inventory    [16]blocks.Stack
	selectedSlot int
	// debug states
	infiniteFuel      bool
	infiniteInventory bool
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

func (t *turtle) move(p coords.Pos) error {
	_, err := t.world.Move(t.pos, p)
	if err != nil {
		return err
	}
	t.pos = p
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
	selection := t.inventory[t.selectedSlot]
	if selection.Count == 0 {
		// NOTE: using equals 0 leaves room for -1 to mean infinite
		return false
	}
	t.inventory[t.selectedSlot] = blocks.Stack{Type: selection.Type, Count: selection.Count - 1}
	toplace := blocks.GetBlock(selection.Type)
	// TODO: placement logic per type should be on block, not on turtle?
	switch selection.Type {
	case blocks.Stairs:
		heading := coords.Pos{t.Heading.Y * -1, t.Heading.X, 0}
		flipped := false
		if t.up() == p {
			if _, upok := t.world.Read(p.Up()); upok {
				flipped = true
			}
		}
		toplace = blocks.BaseBlock{Type: blocks.Stairs, Heading: heading, Flipped: flipped}
	case blocks.CobbleSlab, blocks.BrickSlab:
		flipped := false
		if t.up() == p {
			if _, upok := t.world.Read(p.Up()); upok {
				flipped = true
			}
		}
		toplace = blocks.BaseBlock{Type: selection.Type, Flipped: flipped}
	}
	t.world.Write(p, toplace)
	return true
}

func (t *turtle) Inspect() (blocks.Block, bool) {
	<-t.tick
	b, ok := t.inspect(t.forward())
	t.ack <- true
	return b, ok
}

func (t *turtle) inspect(p coords.Pos) (blocks.Block, bool) {
	return t.world.Read(p)
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
	if t.program == nil {
		return
	}
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
			Type:    blocks.Turtle,
		},
		pos:       p,
		world:     w,
		tick:      tick,
		ack:       ack,
		inventory: [16]blocks.Stack{},
	}
	w.Write(p, t)
	return t
}

func (t *turtle) SetInfiniteInventory() {
	t.infiniteInventory = true
}

func (t *turtle) SetInfiniteFuel() {
	t.infiniteFuel = true
}

// DEBUG statement to set infinite amount of blocks
// assumes selection always at 0
func (t *turtle) SetInventory(bt blocks.Blocktype) {
	if !t.infiniteInventory {
		return
	}
	t.inventory[0] = blocks.Stack{Type: bt, Count: -1}
}
