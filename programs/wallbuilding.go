package programs

import (
    "fmt"
    "github.com/deosjr/TurtleSimulator/blocks"
    "github.com/deosjr/TurtleSimulator/turtle"
)

func debug(t turtle.Turtle, v ...interface{}) {
    fmt.Printf("%#v", t)
    fmt.Println(v...)
}

func sidestepRight(t turtle.Turtle) {
    t.TurnRight()
    t.Forward()
    t.TurnLeft()
}

func sidestepLeft(t turtle.Turtle) {
    t.TurnLeft()
    t.Forward()
    t.TurnRight()
}

func WallCorner() turtle.Program {
    return func(t turtle.Turtle) {
        for i := 0; i < 4; i++ {
            Wallbuildfunc()(t)
            t.Forward()
            t.TurnLeft()
            //debug(t, t.Back())
            t.Back()
        }
        t.TurnRight()
        t.Back()
        t.TurnRight()
        t.Forward()
        t.Forward()
        t.TurnLeft()
        t.Forward()
        t.TurnLeft()
        // right where we started, in place to build the upper part
        Towerfunc(blocks.Stone)(t)
    }
}

func Wallbuildfunc() turtle.Program {
    return func(t turtle.Turtle) {
        t.SetInventory(blocks.Stone)
		if t.Detect() {
			sidestepRight(t)
			t.Forward()
			t.TurnLeft()
		} else {
			t.Forward()
			t.TurnLeft()
			t.Back()
		}
		for i := 0; i < 2; i++ {
			t.Place()
			t.TurnRight()
			t.Place()
			for j := 0; j < 2; j++ {
				sidestepRight(t)
				t.Place()
			}
			t.Up()
			t.Place()
			for j := 0; j < 2; j++ {
				sidestepLeft(t)
				t.Place()
			}
			t.TurnLeft()
			t.Place()
			t.Up()
		}
		t.Place()
		t.Up()
		t.Place()
		t.Down()
		t.TurnRight()
		t.Place()
		t.PlaceUp()

		for i := 0; i < 2; i++ {
			sidestepRight(t)
			t.Place()
			t.PlaceUp()
		}
		t.Down()
		t.TurnRight()
        t.SetInventory(blocks.Stairs)
		t.PlaceUp()
		t.Back()
		t.Back()
		t.TurnRight()
		t.TurnRight()
		t.PlaceUp()
		for i := 0; i < 3; i++ {
			t.Down()
			t.Back()
		}
		t.TurnRight()
		t.Back()
    }
}
