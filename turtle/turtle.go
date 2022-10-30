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
	DigUp() bool
	DigDown() bool
	Inspect() (blocks.Block, bool)
	InspectUp() (blocks.Block, bool)
	InspectDown() (blocks.Block, bool)
	Select(slot int)
	GetSelectedSlot() int
	GetItemCount(slot int) int
	GetItemDetail(slot int) blocks.Blocktype
	GetFuelLevel() int
	Refuel()
	Suck(count int) error
	SuckUp(count int) error
	SuckDown(count int) error

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
	tick    chan struct{}
	ack     chan struct{}
	running bool
	fuel    int
	// inventory management
	inventory    *blocks.Inventory
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
	t.ack <- struct{}{}
}

func (t *turtle) TurnRight() {
	<-t.tick
	t.Heading = coords.Pos{t.Heading.Y * 1, -t.Heading.X, 0}
	t.ack <- struct{}{}
}

func (t *turtle) move(p coords.Pos) error {
	<-t.tick
	defer func() {
		t.ack <- struct{}{}
	}()
	if !t.useFuel() {
		return fmt.Errorf("no fuel!")
	}
	_, err := t.world.Move(t.pos, p)
	if err != nil {
		return err
	}
	t.pos = p
	return nil
}

func (t *turtle) Forward() error {
	return t.move(t.forward())
}

func (t *turtle) Back() error {
	return t.move(t.back())
}

func (t *turtle) Up() error {
	return t.move(t.up())
}

func (t *turtle) Down() error {
	return t.move(t.down())
}

func (t *turtle) Detect() bool {
	return t.detect(t.forward())
}

func (t *turtle) DetectUp() bool {
	return t.detect(t.up())
}

func (t *turtle) DetectDown() bool {
	return t.detect(t.down())
}

func (t *turtle) detect(p coords.Pos) bool {
	_, ok := t.world.Read(p)
	return ok
}

func (t *turtle) Dig() bool {
	return t.dig(t.forward())
}

func (t *turtle) DigUp() bool {
	return t.dig(t.up())
}

func (t *turtle) DigDown() bool {
	return t.dig(t.down())
}

func (t *turtle) dig(p coords.Pos) bool {
	<-t.tick
	defer func() {
		t.ack <- struct{}{}
	}()
	_, ok := t.world.Read(p)
	if !ok {
		return false
	}
	t.world.Delete(p)
	return true
}

func (t *turtle) Place() bool {
	return t.place(t.forward())
}

func (t *turtle) PlaceUp() bool {
	return t.place(t.up())
}

func (t *turtle) PlaceDown() bool {
	return t.place(t.down())
}

func (t *turtle) place(p coords.Pos) bool {
	<-t.tick
	defer func() {
		t.ack <- struct{}{}
	}()
	_, ok := t.world.Read(p)
	if ok {
		return false
	}
	selection := t.inventory.Get(t.selectedSlot)
	if selection.Count == 0 {
		// NOTE: using equals 0 leaves room for -1 to mean infinite
		return false
	}
	t.inventory.Remove(t.selectedSlot, 1)
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
	return t.inspect(t.forward())
}

func (t *turtle) InspectUp() (blocks.Block, bool) {
	return t.inspect(t.up())
}

func (t *turtle) InspectDown() (blocks.Block, bool) {
	return t.inspect(t.down())
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
	t.tick <- struct{}{}
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

// slot should be 0<=slot<16, silent ignore otherwise
// note for generated code that lua indexes from 1 instead
func (t *turtle) Select(slot int) {
	if slot < 0 || slot > 15 {
		return
	}
	t.selectedSlot = slot
}

func (t *turtle) GetSelectedSlot() int {
	return t.selectedSlot
}

func (t *turtle) GetItemCount(slot int) int {
	if slot < 0 || slot > 15 {
		return 0
	}
	return t.inventory.Get(slot).Count
}

func (t *turtle) GetItemDetail(slot int) blocks.Blocktype {
	if slot < 0 || slot > 15 {
		return blocks.Bedrock // TODO invalid?
	}
	return t.inventory.Get(slot).Type
}

func (t *turtle) GetFuelLevel() int {
	return t.fuel
}

// TODO: is this blocking?
func (t *turtle) Refuel() {
	stack := t.inventory.Get(t.selectedSlot)
	// TODO actual fuel calculation other than coal
	if stack.Type != blocks.Coal {
		return
	}
	t.fuel += 8 * stack.Count
	t.inventory.Set(t.selectedSlot, blocks.Stack{})
}

func (t *turtle) useFuel() bool {
	if t.infiniteFuel {
		return true
	}
	if t.fuel < 1 {
		return false
	}
	t.fuel -= 1
	return true
}

func (t *turtle) Suck(count int) error {
	return t.suck(t.forward(), count)
}

func (t *turtle) SuckUp(count int) error {
	return t.suck(t.up(), count)
}

func (t *turtle) SuckDown(count int) error {
	return t.suck(t.down(), count)
}

// TODO: only works on chests rn
func (t *turtle) suck(p coords.Pos, count int) error {
	<-t.tick
	defer func() {
		t.ack <- struct{}{}
	}()
	block, ok := t.world.Read(p)
	if !ok {
		return fmt.Errorf("no block found")
	}
	inv, ok := block.(*blocks.Chest)
	if !ok {
		return fmt.Errorf("no chest found")
	}
	// TODO: assumes selected 0 for both target and dest
	got := inv.Get(0)
	found := got.Count
	if found == 0 {
		return fmt.Errorf("no item found")
	}
	if found <= count {
		inv.Set(0, blocks.Stack{})
		count = found
	} else {
		inv.Remove(0, count)
	}
	if t.inventory.Get(0).Count == 0 {
		t.inventory.Set(0, blocks.Stack{got.Type, found})
	} else {
		t.inventory.Add(0, found)
	}
	return nil
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
	t := &turtle{
		BaseBlock: blocks.BaseBlock{
			Heading: heading,
			Type:    blocks.Turtle,
		},
		pos:       p,
		world:     w,
		tick:      make(chan struct{}, 1),
		ack:       make(chan struct{}, 1),
		inventory: blocks.NewInventory(16),
	}
	w.Write(p, t)
	return t
}

// TODO: perhaps this should be set on inventory struct instead
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
	t.inventory.Set(0, blocks.Stack{Type: bt, Count: -1})
}
