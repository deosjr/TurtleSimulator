package lib

import (
	"github.com/deosjr/TurtleSimulator/turtle"
)

func Turnaround(t turtle.Turtle) {
	t.TurnLeft()
	t.TurnLeft()
}

func UturnLeft(t turtle.Turtle) {
	t.TurnLeft()
	t.Forward()
	t.TurnLeft()
}

func UturnRight(t turtle.Turtle) {
	t.TurnRight()
	t.Forward()
	t.TurnRight()
}
