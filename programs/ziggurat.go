package programs

import (
	"github.com/deosjr/TurtleSimulator/blocks"
	"github.com/deosjr/TurtleSimulator/programs/lib"
	"github.com/deosjr/TurtleSimulator/turtle"
)

//go:generate go run github.com/deosjr/TurtleSimulator/lua

// builds the ziggurat from imgur.com/a/NeywO
// Great Ziggurat by @MCNoodlor

// intended to be built by 4 turtles simultaneously
func Ziggurat(t turtle.Turtle) {
	zigguratPhase1(t)
	zigguratPhase2(t)
	zigguratPhase3(t)
}

func stonesTwoLayers(t turtle.Turtle, n int) {
	for i := 0; i < n-1; i++ {
		t.Back()
		lib.Place(t, blocks.Stone)
		lib.PlaceDown(t, blocks.Stone)
	}
}

func placeStones(t turtle.Turtle, n int) {
	for i := 0; i < n-1; i++ {
		lib.PlaceDown(t, blocks.Stone)
		t.Forward()
	}
	lib.PlaceDown(t, blocks.Stone)
}

// move back while laying front, down AND up
func threeInOne(t turtle.Turtle, n, mod int) {
	for i := mod + 1; i < n+mod+1; i++ {
		t.Back()
		if i%2 == 0 {
			lib.PlaceDown(t, blocks.Planks)
		} else {
			lib.PlaceDown(t, blocks.Stone)
		}
		lib.Place(t, blocks.Planks)
		if i%4 == 0 {
			lib.PlaceUp(t, blocks.Brick)
		} else {
			lib.PlaceUp(t, blocks.BrickSlab)
		}
	}
}

func zigguratPhase1(t turtle.Turtle) {
	t.Forward()
	t.Up()
	t.TurnRight()
	lib.PlaceDown(t, blocks.Stone)
	stonesTwoLayers(t, 37)
	t.TurnRight()
	stonesTwoLayers(t, 8)
	t.Up()
	lib.PlaceDown(t, blocks.Stone)
	t.Up()
	lib.Turnaround(t)
	lib.PlaceDown(t, blocks.Stone)
	lib.PlaceUp(t, blocks.BrickSlab)
	threeInOne(t, 7, 1)
	t.TurnLeft()
	threeInOne(t, 36, 0)
	t.TurnLeft()
	t.Back()
	lib.Place(t, blocks.Planks)
	// now place 4 rows of stone
	t.Up()
	t.TurnRight()
	placeStones(t, 36)
	lib.UturnRight(t)
	placeStones(t, 36)
	lib.UturnLeft(t)
	placeStones(t, 36)
	lib.UturnRight(t)
	placeStones(t, 36)
	lib.Turnaround(t)
	// place wall foundation
	t.Up()
	for i := 0; i < 6; i++ {
		lib.PlaceDown(t, blocks.Stone)
		t.Forward()
		t.Forward()
		t.Forward()
		t.Forward()
	}
	placeStones(t, 9)
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
			lib.PlaceDown(t, blocks.Stone)
		}
	}
	for i := 0; i < n; i++ {
		t.Up()
		lib.PlaceDown(t, blocks.Planks)
		t.Up()
		lib.PlaceDown(t, blocks.Brick)
	}
}

func build3x4block(t turtle.Turtle) {
	lib.PlaceDown(t, blocks.Stone)
	for i := 0; i < 2; i++ {
		t.Back()
		lib.PlaceDown(t, blocks.Stone)
	}
	t.Up()
	t.Forward()
	t.Forward()
	for i := 0; i < 3; i++ {
		lib.PlaceDown(t, blocks.Stone)
		lib.PlaceUp(t, blocks.BrickSlab)
		t.Back()
		lib.Place(t, blocks.Planks)
	}
}

func buildPillars(t turtle.Turtle, n int) {
	for i := 0; i < n; i++ {
		t.Back()
		for j := 0; j < 3; j++ {
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
	for i := 0; i < 5; i++ {
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
		if diffColor > 0 && windowsPlaced == (diffColor-1) {
			windowMat = blocks.Planks
		}
		i := 0
		for !t.DetectDown() {
			lib.Place(t, windowMat)
			sidestepLeft(t)
			lib.Place(t, windowMat)
			sidestepLeft(t)
			lib.Place(t, windowMat)
			t.Down()
			lib.Place(t, windowMat)
			if !stairsPlaced {
				t.TurnLeft()
				lib.PlaceUp(t, blocks.Stairs)
				t.Back()
				t.TurnRight()
			} else {
				sidestepRight(t)
			}
			if i%2 == 0 {
				lib.Place(t, windowMat)
			} else {
				lib.Place(t, blocks.Brick)
			}
			t.TurnRight()
			t.Forward()
			if !stairsPlaced {
				lib.PlaceUp(t, blocks.Stairs)
				stairsPlaced = true
			}
			t.TurnLeft()
			lib.Place(t, windowMat)
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
	for i := 0; i < 3; i++ {
		t.Back()
		lib.Place(t, windowMat)
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
	lib.Turnaround(t)
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
	lib.Place(t, blocks.Stone)
	lib.PlaceUp(t, blocks.CobbleSlab)
}

// starts in the left corner looking at the n-length side
func build3xNfloor(t turtle.Turtle, n int) {
	placeStones(t, n)
	lib.UturnRight(t)
	placeStones(t, n)
	lib.UturnLeft(t)
	placeStones(t, n)
}

func buildStairs(t turtle.Turtle) {
	build3xNfloor(t, 4)
	t.Forward()
	lib.Turnaround(t)
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
	for i := 0; i < 4; i++ {
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
		lib.Place(t, blocks.CobbleSlab)
		t.Up()
		lib.PlaceDown(t, blocks.Stone)
	}

	for i := 0; i < (2*n - 1); i++ {
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
		for i := 0; i < 2; i++ {
			lib.Place(t, windowMat)
			t.Down()
		}
		lib.Place(t, windowMat)
		t.Back()
		lib.Place(t, windowMat)
		for i := 0; i < 2; i++ {
			t.Up()
			lib.Place(t, windowMat)
		}
		for i := 0; i < 3; i++ {
			t.Down()
			lib.PlaceUp(t, windowMat)
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
		lib.Turnaround(t)
		for !t.DetectUp() {
			t.Up()
		}
		return true
	}

	windows(t, bottomFunc2, endFunc2, 0)

	sidestepRight(t)
	t.Up()
	for t.Detect() {
		t.Up()
	}
	for i := 0; i < 3; i++ {
		t.Up()
	}
	t.Forward()
	t.TurnLeft()
	t.Back()
}

func zigguratPhase2(t turtle.Turtle) {
	lib.PlaceUp(t, blocks.Brick)
	for i := 0; i < 3; i++ {
		t.Back()
		lib.PlaceUp(t, blocks.BrickSlab)
		lib.Place(t, blocks.Planks)
	}
	t.Back()
	lib.Place(t, blocks.Planks)

	pillar(t, 2)
	buildPillars(t, 7)
	windows(t, bottomFunc1, endFunc1, 4)
	t.Back()
	t.Back()
	t.TurnLeft()
	for i := 0; i < 10; i++ {
		t.Back()
	}
	buildStairs(t)
	buildUpperWalls(t, 7)
}

func upperpillar(t turtle.Turtle) {
	for _, b := range []blocks.Blocktype{blocks.Planks, blocks.Planks, blocks.Brick, blocks.Stone} {
		t.Up()
		lib.PlaceDown(t, b)
	}
	t.Back()
	lib.Place(t, blocks.Stone)
	t.Down()
	lib.PlaceUp(t, blocks.Stone)
	t.Down()
	lib.PlaceUp(t, blocks.Stairs)
	lib.Turnaround(t)
	t.Forward()
	t.Up()
	lib.PlaceUp(t, blocks.Stone)
	t.Forward()
	lib.PlaceUp(t, blocks.Stone)
	t.Down()
	lib.PlaceUp(t, blocks.Stairs)
}

func upper3x3(t turtle.Turtle) {
	for i := 0; i < 3; i++ {
		if i%2 == 1 {
			lib.PlaceDown(t, blocks.Planks)
		} else {
			lib.PlaceDown(t, blocks.Stone)
		}
		lib.PlaceUp(t, blocks.BrickSlab)
		t.Back()
		lib.Place(t, blocks.Planks)
	}
}

func placeTorchRow(t turtle.Turtle, n int) {
	for i := 0; i < n; i++ {
		t.Back()
		lib.Place(t, blocks.Torch)
		for !t.DetectDown() {
			t.Down()
		}
		for j := 0; j < 3; j++ {
			t.Back()
		}
	}
}

func placeTorchCrossFloor(t turtle.Turtle) {
	t.Back()
	lib.Place(t, blocks.Torch)
	for i := 0; i < 3; i++ {
		t.Back()
	}
	for !t.DetectDown() {
		t.Down()
	}
}

func zigguratPhase3(t turtle.Turtle) {
	t.TurnLeft()
	buildStairs(t)
	buildUpperWalls(t, 5)
	t.TurnLeft()
	buildStairs(t)
	findBottomNextCorner(t)
	t.Up()
	lib.PlaceDown(t, blocks.Stone)
	t.Up()
	for i := 0; i < 11; i++ {
		lib.PlaceDown(t, blocks.Stone)
		lib.PlaceUp(t, blocks.Stone)
		t.Back()
		lib.Place(t, blocks.Stone)
	}
	for i := 0; i < 3; i++ {
		lib.PlaceDown(t, blocks.Stone)
		t.Up()
	}
	lib.Turnaround(t)
	for i := 0; i < 11; i++ {
		lib.PlaceDown(t, blocks.Stone)
		if i%4 == 0 {
			lib.PlaceUp(t, blocks.Brick)
		} else {
			lib.PlaceUp(t, blocks.BrickSlab)
		}
		t.Back()
		lib.Place(t, blocks.Planks)
	}
	lib.PlaceDown(t, blocks.Stone)
	lib.PlaceUp(t, blocks.BrickSlab)
	t.TurnRight()
	t.Back()
	lib.Place(t, blocks.Planks)
	t.Up()
	t.Back()
	t.Back()
	t.TurnLeft()
	build3xNfloor(t, 9)
	t.Forward()
	t.TurnLeft()
	build3xNfloor(t, 5)
	t.TurnLeft()
	t.Forward()
	lib.Place(t, blocks.BrickSlab)
	t.Back()
	lib.Place(t, blocks.Brick)

	t.TurnLeft()
	t.Up()
	for i := 0; i < 5; i++ {
		t.Forward()
	}
	t.TurnLeft()
	upperpillar(t)
	t.Forward()
	t.Down()
	t.Down()
	lib.Turnaround(t)
	upperpillar(t)
	t.Forward()
	t.Down()
	t.Down()
	t.TurnLeft()
	upperpillar(t)
	t.TurnRight()
	t.Forward()
	for i := 0; i < 4; i++ {
		t.Up()
	}
	t.Back()
	t.TurnLeft()

	upper3x3(t)
	t.TurnLeft()
	lib.PlaceDown(t, blocks.Planks)
	lib.PlaceUp(t, blocks.Brick)
	t.Back()
	lib.Place(t, blocks.Planks)
	upper3x3(t)
	lib.PlaceDown(t, blocks.Planks)
	lib.PlaceUp(t, blocks.Brick)
	t.Back()
	lib.Place(t, blocks.Planks)
	upper3x3(t)
	lib.PlaceDown(t, blocks.Planks)
	lib.PlaceUp(t, blocks.Brick)
	t.TurnLeft()
	t.Back()
	lib.Place(t, blocks.Planks)
	t.Up()
	t.TurnRight()
	build3xNfloor(t, 8)
	lib.Turnaround(t)
	for i := 0; i < 4; i++ {
		t.Forward()
	}
	sidestepLeft(t)
	lib.PlaceDown(t, blocks.Planks)
	for i := 0; i < 3; i++ {
		t.Forward()
		lib.PlaceDown(t, blocks.Planks)
	}
	lib.Turnaround(t)
	t.Up()
	lib.PlaceDown(t, blocks.Brick)
	for i := 0; i < 3; i++ {
		t.Forward()
		lib.PlaceDown(t, blocks.BrickSlab)
	}

	for i := 0; i < 3; i++ {
		t.Back()
	}
	t.TurnRight()
	placeTorchCrossFloor(t)
	t.TurnRight()
	placeTorchRow(t, 3)
	t.TurnRight()
	placeTorchRow(t, 5)
	t.TurnRight()
	placeTorchRow(t, 7)
	t.TurnRight()
	placeTorchRow(t, 8)
	placeTorchCrossFloor(t)
	t.TurnLeft()
	placeTorchRow(t, 1)
	t.TurnLeft()
	placeTorchRow(t, 5)
	// in case we pass over entry stairs
	placeTorchCrossFloor(t)
	placeTorchRow(t, 3)
	t.TurnRight()
	t.Back()
	lib.Place(t, blocks.Torch)
	for i := 0; i < 5; i++ {
		t.Down()
	}
}
