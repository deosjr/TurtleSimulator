package programs

import (
    "github.com/deosjr/TurtleSimulator/blocks"
    "github.com/deosjr/TurtleSimulator/turtle"
)

func turnaround(t turtle.Turtle) {
    t.TurnLeft()
    t.TurnLeft()
}

func placeFringe(t turtle.Turtle) {
    endReached := false
	for !endReached {
		t.Back()
		place(t, blocks.Stairs)
		t.Back()
		turnaround(t)
		endReached = !place(t, blocks.Brick)
		t.Back()
		place(t, blocks.Stairs)
		t.Up()
		placeDown(t, blocks.Slab)
		t.Forward()
		if !endReached {
			t.Forward()
			t.Forward()
			turnaround(t)
			place(t, blocks.Slab)
			t.Down()
		}
	}
}

func placeFloorRow(t turtle.Turtle) {
	for !t.Detect() {
		placeDown(t, blocks.Stone)
		t.Forward()
	}
	placeDown(t, blocks.Stone)
}

func placeFloor(t turtle.Turtle) {
    t.Down()
	turnaround(t)
	placeFloorRow(t)
	t.TurnRight()
	t.Forward()
	t.TurnRight()
	placeFloorRow(t)
	t.TurnLeft()
	t.Forward()
	t.TurnLeft()
	placeFloorRow(t)
}

func Upperwallfunc() turtle.Program {
    return func(t turtle.Turtle) {
        for i:=0;i<8;i++ {
            t.Up()
        }
        sidestepRight(t)
        t.Forward()
        t.TurnLeft()
        for i:=0;i<2;i++ {
	        t.Down()
	        placeDown(t, blocks.Stone)
	        t.Back()
            for !t.DetectDown() {
		        place(t, blocks.Log)
		        placeDown(t, blocks.Log)
		        t.Up()
		        placeDown(t, blocks.Log)
		        t.Back()
		        t.Down()
		        placeDown(t, blocks.Stone)
		        t.Back()
            }
	        t.Dig()
	        t.Forward()
	        turnaround(t)
	        place(t, blocks.Log)
	        t.Up()
	        placeDown(t, blocks.Log)
	        t.Forward()
	        sidestepLeft(t)
	        t.Down()
	        t.Dig()
	        place(t, blocks.Stone)
	        t.Down()
	        t.Dig()
	        sidestepLeft(t)
	        t.Dig()
	        t.Up()
	        t.Dig()
	        place(t, blocks.Stone)
	        sidestepLeft(t)
	        t.Dig()
	        place(t, blocks.Stone)
	        t.Down()
	        t.Dig()
	        t.Up()
	        t.Up()
	        sidestepLeft(t)
        }
        placeFringe(t)
        sidestepLeft(t)
        placeFloor(t)
        sidestepRight(t)
        placeFringe(t)

        t.TurnLeft()
        for i:=0;i<4;i++ {
	        t.Back()
        }
        t.TurnRight()
        /*
        for i=1,(x-1) do
	        turtle.back()
        end
        */
        t.TurnLeft()
        t.Back()
        for i:=0;i<9;i++ {
	        t.Down()
        }
    }
}
