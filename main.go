package main

import (
    "github.com/deosjr/TurtleSimulator/blocks"
    "github.com/deosjr/TurtleSimulator/coords"
    "github.com/deosjr/TurtleSimulator/programs"
    "github.com/deosjr/TurtleSimulator/turtle"
)


// TODO: a game tick in minecraft is 1/20 second
// from dan200 computercraft java code:
// each animation takes 8 ticks to complete unless otherwise specified.

func main() {
	tick := make(chan bool, 1)
	tack := make(chan bool, 1)
    w := turtle.NewWorld(5, tick, tack)
	//w.grid[pos{0, 3, 0}] = block{}
	//w.grid[pos{2, 2, 0}] = block{}
	for y := -10; y < 10; y++ {
		for x := -10; x < 10; x++ {
			w.Grid[coords.Pos{x, y, -1}] = blocks.GetBlock(blocks.Grass)
		}
	}
	t := turtle.NewTurtle(coords.Pos{0, 0, 0}, w, tick, tack)
    w.Turtles = []turtle.Turtle{t}
	// wall building program
	t.SetProgram(programs.WallCorner())

    vis := NewRaytracer(false, false)
    //vis := ascii{}
    //Visualise(vis, w)
    VisualiseEndState(vis, w)
}
