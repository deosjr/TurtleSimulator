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
    w := turtle.NewWorld(5)
	for y := -10; y < 10; y++ {
		for x := -10; x < 10; x++ {
            w.Write(coords.Pos{x, y, -1}, blocks.GetBlock(blocks.Grass))
		}
	}
	t1 := turtle.NewTurtle(coords.Pos{0, 0, 0}, w, coords.North)
	t2 := turtle.NewTurtle(coords.Pos{5, 5, 0}, w, coords.West)
    w.Turtles = []turtle.Turtle{t1, t2}
	// wall building program
	t1.SetProgram(programs.Walls(9, "L9", "L9", "L9"))
	t2.SetProgram(programs.Walls(7, "R7", "R7", "R7"))

    vis := NewRaytracer(false, false)
    //vis := ascii{}
    //Visualise(vis, w)
    VisualiseEndState(vis, w)
}
