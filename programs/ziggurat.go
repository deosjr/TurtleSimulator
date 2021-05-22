package programs

import (
    "github.com/deosjr/TurtleSimulator/blocks"
    "github.com/deosjr/TurtleSimulator/turtle"
)

// builds the ziggurat from imgur.com/a/NeywO
// Great Ziggurat by @MCNoodlor

// intended to be built by 4 turtles simultaneously
func Ziggurat() turtle.Program {
    return func(t turtle.Turtle) {
        ZigguratPhase1()(t)
        ZigguratPhase2()(t)
    }
}

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

// go down until detectdown, build stone back up until planks in front,
// then top it off with n stacks of planks + brick
func pillar(t turtle.Turtle, n int) {
    for !t.DetectDown() {
        t.Down()
    }
	planksDetected := false
	for !planksDetected {
		data, ok := t.Inspect()
		if ok && data.GetType() == blocks.Planks {
			planksDetected = true
		} else {
			t.Up()
			placeDown(t, blocks.Stone)
		}
	}
	for i:=0; i<n; i++ {
		t.Up()
		placeDown(t, blocks.Planks)
		t.Up()
		placeDown(t, blocks.Brick)
	}
}

func build3x4block(t turtle.Turtle) {
	placeDown(t, blocks.Stone)
	for i:=0; i<2; i++ {
		t.Back()
		placeDown(t, blocks.Stone)
	}
	t.Up()
	t.Forward()
	t.Forward()
	for i:=0; i<3; i++ {
		placeDown(t, blocks.Stone)
		placeUp(t, blocks.BrickSlab)
		t.Back()
		place(t, blocks.Planks)
	}
}

func buildPillars(t turtle.Turtle, n int) {
	for i:=0; i<n; i++ {
		t.Back()
	    for j:=0; j<3; j++ {
			t.Down()
		}
		build3x4block(t)
		if i == (n-2) || i == (n-1) {
			pillar(t, 1)
		} else {
			pillar(t, 2)
		}
	}
}

type bottomFunc func(turtle.Turtle, blocks.Blocktype)
type endFunc func(t turtle.Turtle) bool

// func when bottom of window reached
// func when end of windows reached
// num of window that should be different block, if any
func windows(t turtle.Turtle, bottomFunc bottomFunc, endFunc endFunc, diffColor int) {
	t.Forward()
	sidestepRight(t)
    for i:=0; i<5; i++ {
		t.Down()
	}
	t.TurnRight()
	t.Back()
	// build 6 stones, three left, down, then three right
	// until we hit the ground
	// every second time, place brick in the middle of lower layer
	// fourth window is made of planks
	moreWindows := true
	windowsPlaced := 0
	for moreWindows {
		stairsPlaced := false
		windowMat := blocks.Stone
		if diffColor > 0 && windowsPlaced == (diffColor - 1) {
			windowMat = blocks.Planks
		}
		i := 0
		for !t.DetectDown() {
			place(t, windowMat)
			sidestepLeft(t)
			place(t, windowMat)
			sidestepLeft(t)
			place(t, windowMat)
			t.Down()
			place(t, windowMat)
			if !stairsPlaced {
				t.TurnLeft()
				placeUp(t, blocks.Stairs)
				t.Back()
				t.TurnRight()
			} else {
				sidestepRight(t)
			}
			if i%2==0 {
				place(t, windowMat)
			} else {
				place(t, blocks.Brick)
			}
			t.TurnRight()
		    t.Forward()
			if !stairsPlaced {
				placeUp(t, blocks.Stairs)
				stairsPlaced = true
			}
			t.TurnLeft()
			place(t, windowMat)
			t.Down()
			i = (i + 1) % 2
		}

		bottomFunc(t, windowMat)
		windowsPlaced = windowsPlaced + 1
		moreWindows = endFunc(t)
	}
}

func bottomFunc1(t turtle.Turtle, windowMat blocks.Blocktype) {
	// place the last three blocks and find next window
	t.Forward()
	t.TurnRight()
    for i:=0; i<3; i++ {
		t.Back()
		place(t, windowMat)
	}
	t.Back()
	t.TurnRight()
}

func endFunc1(t turtle.Turtle) bool {
// turtle is looking through next window if any
	if t.Detect() {
		return false
	}
	t.Forward()
	turnaround(t)
    for !t.DetectUp() {
        t.Up()
    }
	return true
}

// start stairs
func backupback(t turtle.Turtle) {
	t.Back()
	t.Up()
	t.Back()
}

func stairpiece(t turtle.Turtle) {
	place(t, blocks.Stone)
	placeUp(t, blocks.CobbleSlab)
}

// starts in the left corner looking at the n-length side
func build3xNfloor(t turtle.Turtle, n int) {
	placeStones(t, n)
	uturnRight(t)
	placeStones(t, n)
	uturnLeft(t)
	placeStones(t, n)
}

func buildStairs(t turtle.Turtle) {
	build3xNfloor(t, 4)
	t.Forward()
	turnaround(t)
	t.Down()

	stairsDone := false
	for !stairsDone {
		stairpiece(t)
		sidestepRight(t)
		stairpiece(t)
		sidestepRight(t)
		stairpiece(t)
		backupback(t)

		stairpiece(t)
		sidestepLeft(t)
		stairpiece(t)
		sidestepLeft(t)
		stairpiece(t)
		backupback(t)

		t.TurnLeft()
        data, ok := t.Inspect()
		if ok && data.GetType() == blocks.Planks {
			backupback(t)
			sidestepRight(t)
			placeStones(t, 3)
			stairsDone = true
		} else {
			t.TurnRight()
		}
	}
}
// end stairs

func findBottomNextCorner(t turtle.Turtle) {
	t.TurnLeft()
    for i:=0; i<4; i++ {
		t.Forward()
	}
	t.TurnRight()
	for t.Detect() {
		t.Down()
	}
	t.Up()
	t.Up()
}

// n is number of wall pieces
func buildUpperWalls(t turtle.Turtle, n int) {
	findBottomNextCorner(t)

	slabblock := func() {
		t.Back()
		place(t, blocks.CobbleSlab)
		t.Up()
		placeDown(t, blocks.Stone)
	}

	for i:=0;i<(2*n-1);i++ {
		slabblock()
		t.Back()
	}
	slabblock()

	t.Up()
    for !t.DetectDown() {
        t.Forward()
    }
	// start wall up
	buildPillars(t, n)

	bottomFunc2 := func(t turtle.Turtle, windowMat blocks.Blocktype) {
		t.Forward()
		t.TurnLeft()
		t.Forward()
        for i:=0; i<2; i++ {
			place(t, windowMat)
			t.Down()
		}
		place(t, windowMat)
		t.Back()
		place(t, windowMat)
        for i:=0; i<2; i++ {
			t.Up()
			place(t, windowMat)
		}
        for i:=0; i<3; i++ {
			t.Down()
			placeUp(t, windowMat)
		}
		t.Forward()
	}

	endFunc2 := func(t turtle.Turtle) bool {
		// turtle is looking at previous wall if we are done
		if t.Detect() {
			// hit a corner piece going up
			return false
		}
		t.Forward()
		if t.Detect() {
			// hit a corner piece at equal height
			t.Back()
			return false
		}
		t.Forward()
		t.Up()
		t.Forward()
		t.TurnLeft()
		t.Forward()
		turnaround(t)
        for !t.DetectUp() {
			t.Up()
        }
		return true
	}

	windows(t, bottomFunc2, endFunc2, 0)

	sidestepRight(t)
    for t.Detect() {
		t.Up()
    }
    for i:=0; i<3; i++ {
		t.Up()
	}
	t.Forward()
	t.TurnLeft()
	t.Back()
}

func ZigguratPhase2() turtle.Program {
    return func(t turtle.Turtle) {
        placeUp(t, blocks.Brick)
	    for i:=0; i<3; i++ {
		    t.Back()
		    placeUp(t, blocks.BrickSlab)
		    place(t, blocks.Planks)
	    }
	    t.Back()
	    place(t, blocks.Planks)

	    pillar(t, 2)
	    buildPillars(t, 7)
	    windows(t, bottomFunc1, endFunc1, 4)
	    t.Back()
	    t.Back()
	    t.TurnLeft()
        for i:=0; i<10; i++ {
            t.Back()
        }
	    buildStairs(t)
	    buildUpperWalls(t, 7)
    }
}
