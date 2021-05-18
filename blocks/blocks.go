package blocks

import (
    "github.com/deosjr/TurtleSimulator/coords"
)

type Block interface{
    GetHeading() coords.Pos
    GetType() Blocktype
}

type BaseBlock struct{
    Heading coords.Pos
    Flipped bool
    Type Blocktype
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
    Turtle
    Stone
    Grass
    Stairs
    Log
    Planks
    Brick
    CobbleSlab
    BrickSlab
)

func GetBlock(t Blocktype) Block {
    return BaseBlock{Type:t}
}
