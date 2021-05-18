package programs

import (
    "github.com/deosjr/TurtleSimulator/blocks"
    "github.com/deosjr/TurtleSimulator/turtle"
)

// builds the ziggurat from imgur.com/a/NeywO
// Great Ziggurat by @MCNoodlor

// intended to be built by 4 turtles simultaneously

func stonesTwoLayers(t turtle.Turtle, n int) {
    for i:=0; i<n-1; i++ {
        t.Back()
        place(t, blocks.Stone)
        placeDown(t, blocks.Stone)
    }
}

func placeStones(t turtle.Turtle, n int) {
    for i:=0; i<n-1; i++ {
		placeDown(t, blocks.Stone)
		t.Forward()
    }
	placeDown(t, blocks.Stone)
}

// move back while laying front, down AND up
func threeInOne(t turtle.Turtle, n, mod int) {
    for i:=mod; i<n+mod; i++ {
		t.Back()
		if i%2==0 {
			placeDown(t, blocks.Planks)
		} else {
			placeDown(t, blocks.Stone)
		}
		place(t, blocks.Planks)
		if i%4==0 {
			placeUp(t, blocks.Brick)
		} else {
			placeUp(t, blocks.BrickSlab)
		}
	}
}

func ZigguratPhase1() turtle.Program {
    return func(t turtle.Turtle) {
        t.Forward()
	    t.Up()
	    t.TurnRight()
	    placeDown(t, blocks.Stone)
	    stonesTwoLayers(t, 37)
	    t.TurnRight()
	    stonesTwoLayers(t, 8)
	    t.Up()
	    placeDown(t, blocks.Stone)
	    t.Up()
	    turnaround(t)
	    placeDown(t, blocks.Stone)
	    placeUp(t, blocks.BrickSlab)
	    threeInOne(t, 7, 1)
	    t.TurnLeft()
	    threeInOne(t, 36, 0)
	    t.TurnLeft()
	    t.Back()
	    place(t, blocks.Planks)
	    // now place 4 rows of stone
	    t.Up()
	    t.TurnRight()
	    placeStones(t, 36)
	    uturnRight(t)
	    placeStones(t, 36)
	    uturnLeft(t)
	    placeStones(t, 36)
	    uturnRight(t)
	    placeStones(t, 36)
	    turnaround(t)
	    // place wall foundation
	    t.Up()
	    for i:=0; i<6; i++ {
		    placeDown(t, blocks.Stone)
		    t.Forward()
		    t.Forward()
		    t.Forward()
		    t.Forward()
	    }
	    placeStones(t, 9)
    }
}
