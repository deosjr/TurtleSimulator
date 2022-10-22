package programs

import (
	"github.com/deosjr/TurtleSimulator/blocks"
	"github.com/deosjr/TurtleSimulator/programs/lib"
	"github.com/deosjr/TurtleSimulator/turtle"
)

// upper part of a tower
// to be placed on top of a wall corner

func Towerfunc(dummy blocks.Blocktype) turtle.Program {
	return func(t turtle.Turtle) {
		y := 0
		for t.Detect() {
			t.Up()
			y++
		}
		for i := 0; i < 4; i++ {
			buildTowerFace(t, i, dummy)
		}
		for i := 0; i < y; i++ {
			t.Down()
		}
	}
}

func buildTowerFace(t turtle.Turtle, faceNum int, dummy blocks.Blocktype) {
	t.Forward()
	t.TurnRight()
	t.Up()
	lib.PlaceDown(t, blocks.Log)
	t.Forward()
	lib.PlaceDown(t, blocks.Stone)
	t.Forward()
	lib.PlaceDown(t, blocks.Log)
	t.Forward()
	lib.PlaceDown(t, blocks.Stone)
	usedummy := !t.Detect()
	if usedummy {
		lib.Place(t, dummy)
	}
	t.Back()
	lib.Place(t, blocks.Log)
	t.Up()
	lib.PlaceDown(t, blocks.Log)
	t.Back()
	t.Down()
	t.Back()
	lib.Place(t, blocks.Log)
	t.Up()
	lib.PlaceDown(t, blocks.Log)

	t.TurnLeft()
	t.Back()
	for i, bt := range []blocks.Blocktype{blocks.Brick, blocks.Planks, blocks.Planks, blocks.Planks, blocks.Brick, blocks.BrickSlab} {
		if i > 0 {
			t.Up()
		}
		lib.Place(t, bt)
	}
	sidestepRight(t)
	t.Down()
	t.Forward()
	t.TurnLeft()
	t.Back()
	lib.Place(t, blocks.Stairs)
	sidestepLeft(t)
	t.Forward()
	t.TurnRight()
	for _, bt := range []blocks.Blocktype{blocks.Planks, blocks.Brick, blocks.Brick, blocks.Brick} {
		t.Down()
		lib.Place(t, bt)
	}
	sidestepRight(t)
	lib.Place(t, blocks.BrickSlab)
	for i := 0; i < 3; i++ {
		t.Up()
	}
	lib.Place(t, blocks.Planks)
	t.Up()
	t.Forward()
	t.TurnRight()
	t.Forward()
	lib.PlaceDown(t, blocks.Planks)
	t.Back()
	lib.Place(t, blocks.Stairs)
	t.TurnLeft()
	t.Back()
	lib.Place(t, blocks.BrickSlab)
	sidestepRight(t)
	t.Down()
	for i := 0; i < 3; i++ {
		t.Down()
		lib.Place(t, blocks.Brick)
	}

	t.Up()
	sidestepLeft(t)
	t.Forward()
	lib.PlaceUp(t, blocks.BrickSlab)
	t.Forward()
	t.Down()
	lib.PlaceDown(t, blocks.Stone)
	t.TurnLeft()
	t.Forward()
	lib.PlaceDown(t, blocks.Stone)
	t.Up()
	t.Up()
	lib.PlaceUp(t, blocks.Stone)
	t.Back()
	lib.PlaceUp(t, blocks.Stone)
	t.TurnLeft()

	if faceNum == 3 {
		t.Back()
		lib.PlaceUp(t, blocks.Stone)
		t.Down()
		t.Down()
		lib.PlaceDown(t, blocks.Stone)
		t.Up()
		t.Forward()
	} else {
		t.Down()
	}

	t.Forward()
	t.Forward()
	t.TurnRight()
	for i := 0; i < 3; i++ {
		t.Back()
		t.Down()
		if usedummy && i == 1 {
			// clean up the dummy block on the way back
			t.TurnRight()
			t.Dig()
			t.TurnLeft()
		}
	}
	sidestepRight(t)
}
