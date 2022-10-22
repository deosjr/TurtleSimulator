package scene

import (
	"fmt"

	"github.com/deosjr/TurtleSimulator/coords"
	"github.com/deosjr/TurtleSimulator/programs"
	"github.com/deosjr/TurtleSimulator/turtle"
)

func Fort() *turtle.World {
	// TODO: perhaps set camera angles from here as well?
	w := turtle.NewWorld(5)

	n := 10

	t1 := turtle.NewTurtle(coords.Pos{0, 0, 0}, w, coords.North)
	t2 := turtle.NewTurtle(coords.Pos{n*4 + 1, 1, 0}, w, coords.West)
	t3 := turtle.NewTurtle(coords.Pos{n * 4, n*4 + 2, 0}, w, coords.South)
	t4 := turtle.NewTurtle(coords.Pos{-1, n*4 + 1, 0}, w, coords.East)
	t5 := turtle.NewTurtle(coords.Pos{5, 5, 0}, w, coords.West)
	w.Turtles = []turtle.Turtle{t1, t2, t3, t4, t5}
	// wall building program
	t1.SetProgram(programs.Outerwall(n))
	t2.SetProgram(programs.Outerwall(n))
	t3.SetProgram(programs.Outerwall(n))
	t4.SetProgram(programs.Outerwall(n))
	right := fmt.Sprintf("R%d", n-2)
	t5.SetProgram(programs.Walls(n-2, right, right, right))

	t1.SetInfiniteInventory()
	t2.SetInfiniteInventory()
	t3.SetInfiniteInventory()
	t4.SetInfiniteInventory()
	t5.SetInfiniteInventory()

	return w
}
