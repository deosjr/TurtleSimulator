package turtle

import (
    "github.com/deosjr/TurtleSimulator/blocks"
    "github.com/deosjr/TurtleSimulator/coords"
)

// lives here until i promote it to yet another package

// convention: world minimal coord is 0,0,0 and max dim,dim,dim
type World struct {
	Grid map[coords.Pos]blocks.Block
	Dim  int
    Turtles []Turtle
    Tick chan bool
    Tack chan bool
}

func NewWorld(dim int, tick, tack chan bool) *World {
    return &World{
        Grid: map[coords.Pos]blocks.Block{},
        Dim: dim,
        Tick: tick,
        Tack: tack,
    }
}
