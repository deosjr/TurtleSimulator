package blocks

import (
	"github.com/deosjr/TurtleSimulator/coords"
)

type Block interface {
	GetHeading() coords.Pos
	GetType() Blocktype
}

type BaseBlock struct {
	Heading coords.Pos
	Flipped bool
	Type    Blocktype
}

func (b BaseBlock) GetHeading() coords.Pos {
	return b.Heading
}

func (b BaseBlock) GetType() Blocktype {
	return b.Type
}

type Blocktype int

const (
	Bedrock Blocktype = iota
	Stone
	Grass
	Stairs
	Log
	Planks
	Brick
	CobbleSlab
	BrickSlab
	Torch
	// TODO: items!
	Coal
	// TODO: complex blocks with state separate?
	ChestType
	Turtle
)

func GetBlock(t Blocktype) Block {
	return BaseBlock{Type: t}
}

type Stack struct {
	Type  Blocktype
	Count int
}

type Inventory struct {
	NumSlots int
	Slots    []Stack
}

func NewInventory(size int) *Inventory {
	return &Inventory{NumSlots: size, Slots: make([]Stack, size)}
}

func (i *Inventory) Get(slot int) Stack {
	if slot < 0 || slot > i.NumSlots {
		return Stack{}
	}
	return i.Slots[slot]
}

func (i *Inventory) Add(slot, n int) {
	if slot < 0 || slot > i.NumSlots {
		return
	}
	stack := i.Slots[slot]
	i.Slots[slot] = Stack{Type: stack.Type, Count: stack.Count + n}
}

func (i *Inventory) Remove(slot, n int) {
	if slot < 0 || slot > i.NumSlots {
		return
	}
	stack := i.Slots[slot]
	i.Slots[slot] = Stack{Type: stack.Type, Count: stack.Count - n}
}

func (i *Inventory) Set(slot int, stack Stack) {
	if slot < 0 || slot > i.NumSlots {
		return
	}
	i.Slots[slot] = stack
}

type Chest struct {
	BaseBlock
	Inventory
}
