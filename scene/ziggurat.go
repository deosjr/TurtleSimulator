package scene

import (
    "github.com/deosjr/TurtleSimulator/coords"
    "github.com/deosjr/TurtleSimulator/programs"
    "github.com/deosjr/TurtleSimulator/turtle"
)

func Ziggurat() *turtle.World {
    w := turtle.NewWorld(5)

	t1 := turtle.NewTurtle(coords.Pos{20, 0, 0}, w, coords.North)
    w.Turtles = []turtle.Turtle{t1}
    t1.SetProgram(programs.ZigguratPhase1())

    return w
}
