package scene

import (
	"github.com/deosjr/TurtleSimulator/coords"
	"github.com/deosjr/TurtleSimulator/programs"
	"github.com/deosjr/TurtleSimulator/turtle"
)

func Ziggurat() *turtle.World {
	w := turtle.NewWorld(5)

	t1 := turtle.NewTurtle(coords.Pos{37, 0, 0}, w, coords.North)
	t2 := turtle.NewTurtle(coords.Pos{46, 37, 0}, w, coords.West)
	t3 := turtle.NewTurtle(coords.Pos{9, 46, 0}, w, coords.South)
	t4 := turtle.NewTurtle(coords.Pos{0, 9, 0}, w, coords.East)
	w.Turtles = []turtle.Turtle{t1, t2, t3, t4}
	t1.SetProgram(programs.Ziggurat)
	t2.SetProgram(programs.Ziggurat)
	t3.SetProgram(programs.Ziggurat)
	t4.SetProgram(programs.Ziggurat)

	return w
}
