package programs

import (
	"github.com/deosjr/TurtleSimulator/blocks"
	"github.com/deosjr/TurtleSimulator/programs/lib"
	"github.com/deosjr/TurtleSimulator/turtle"
)

func placeFringe(t turtle.Turtle) {
	endReached := false
	for !endReached {
		t.Back()
		lib.Place(t, blocks.Stairs)
		t.Back()
		lib.Turnaround(t)
		endReached = !lib.Place(t, blocks.Brick)
		t.Back()
		lib.Place(t, blocks.Stairs)
		t.Up()
		lib.PlaceDown(t, blocks.BrickSlab)
		t.Forward()
		if !endReached {
			t.Forward()
			t.Forward()
			lib.Turnaround(t)
			lib.Place(t, blocks.BrickSlab)
			t.Down()
		}
	}
}

func placeFloorRow(t turtle.Turtle) {
	for !t.Detect() {
		lib.PlaceDown(t, blocks.Stone)
		t.Forward()
	}
	lib.PlaceDown(t, blocks.Stone)
}

func placeFloor(t turtle.Turtle) {
	t.Down()
	lib.Turnaround(t)
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

func Upperwallfunc(t turtle.Turtle) {
	for i := 0; i < 8; i++ {
		t.Up()
	}
	sidestepRight(t)
	t.Forward()
	t.TurnLeft()
	for i := 0; i < 2; i++ {
		t.Down()
		lib.PlaceDown(t, blocks.Stone)
		t.Back()
		for !t.DetectDown() {
			lib.Place(t, blocks.Log)
			lib.PlaceDown(t, blocks.Log)
			t.Up()
			lib.PlaceDown(t, blocks.Log)
			t.Back()
			t.Down()
			lib.PlaceDown(t, blocks.Stone)
			t.Back()
		}
		t.Dig()
		t.Forward()
		lib.Turnaround(t)
		lib.Place(t, blocks.Log)
		t.Up()
		lib.PlaceDown(t, blocks.Log)
		t.Forward()
		sidestepLeft(t)
		t.Down()
		t.Dig()
		lib.Place(t, blocks.Stone)
		t.Down()
		t.Dig()
		sidestepLeft(t)
		t.Dig()
		t.Up()
		t.Dig()
		lib.Place(t, blocks.Stone)
		sidestepLeft(t)
		t.Dig()
		lib.Place(t, blocks.Stone)
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
	for i := 0; i < 4; i++ {
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
	for i := 0; i < 9; i++ {
		t.Down()
	}
}
