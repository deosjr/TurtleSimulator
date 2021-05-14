package programs

import (
    "github.com/deosjr/TurtleSimulator/blocks"
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

// helper funcs, should later on also include selecting from inv if possible
func place(t turtle.Turtle, bt blocks.Blocktype) {
    t.SetInventory(bt)
    t.Place()
}
func placeUp(t turtle.Turtle, bt blocks.Blocktype) {
    t.SetInventory(bt)
    t.PlaceUp()
}
func placeDown(t turtle.Turtle, bt blocks.Blocktype) {
    t.SetInventory(bt)
    t.PlaceDown()
}

func buildTowerFace(t turtle.Turtle, faceNum int, dummy blocks.Blocktype) {
    t.Forward()
    t.TurnRight()
    t.Up()
    placeDown(t, blocks.Log)
    t.Forward()
    placeDown(t, blocks.Stone)
    t.Forward()
    placeDown(t, blocks.Log)
    t.Forward()
    placeDown(t, blocks.Stone)
    usedummy := !t.Detect()
    if usedummy {
        place(t, dummy)
    }
    t.Back()
    place(t, blocks.Log)
    t.Up()
    placeDown(t, blocks.Log)
    t.Back()
    t.Down()
    t.Back()
    place(t, blocks.Log)
    t.Up()
    placeDown(t, blocks.Log)

    t.TurnLeft()
    t.Back()
    for i, bt := range []blocks.Blocktype{blocks.Brick, blocks.Planks, blocks.Planks, blocks.Planks, blocks.Brick, blocks.Slab} {
        if i > 0 {
            t.Up()
        }
        place(t, bt)
    }
	sidestepRight(t)
	t.Down()
	t.Forward()
	t.TurnLeft()
	t.Back()
	place(t, blocks.Stairs)
	sidestepLeft(t)
	t.Forward()
	t.TurnRight()
    for _, bt := range []blocks.Blocktype{blocks.Planks, blocks.Brick, blocks.Brick, blocks.Brick} {
        t.Down()
        place(t, bt)
    }
	sidestepRight(t)
	place(t, blocks.Slab)
    for i:=0;i<3;i++ {
        t.Up()
    }
    place(t, blocks.Planks)
	t.Up()
	t.Forward()
	t.TurnRight()
	t.Forward()
	placeDown(t, blocks.Planks)
	t.Back()
	place(t, blocks.Stairs)
	t.TurnLeft()
	t.Back()
	place(t, blocks.Slab)
	sidestepRight(t)
	t.Down()
    for i:=0;i<3;i++ {
		t.Down()
		place(t, blocks.Brick)
	}

    t.Up()
	sidestepLeft(t)
	t.Forward()
	placeUp(t, blocks.Slab)
	t.Forward()
	t.Down()
	placeDown(t, blocks.Stone)
	t.TurnLeft()
	t.Forward()
	placeDown(t, blocks.Stone)
	t.Up()
	t.Up()
	placeUp(t, blocks.Stone)
	t.Back()
	placeUp(t, blocks.Stone)
	t.TurnLeft()

    if faceNum == 3 {
		t.Back()
		placeUp(t, blocks.Stone)
		t.Down()
		t.Down()
		placeDown(t, blocks.Stone)
		t.Up()
		t.Forward()
	} else {
		t.Down()
	}

    t.Forward()
	t.Forward()
	t.TurnRight()
    for i:=0;i<3;i++ {
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
