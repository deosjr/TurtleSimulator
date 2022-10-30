package turtle

import (
	"testing"

	"github.com/deosjr/TurtleSimulator/blocks"
	"github.com/deosjr/TurtleSimulator/coords"
)

func TestMovement(t *testing.T) {
	w := NewWorld(5)
	turtle := NewTurtle(coords.Pos{0, 0, 0}, w, coords.North).(*turtle)
	turtle.SetInfiniteFuel()
	turtle.SetProgram(func(tt Turtle) {
		tt.Forward()
		tt.TurnRight()
		tt.Forward()
	})
	// TODO: does this mean the world only knows about turtle
	// positions once they start moving?
	w.Turtles = []Turtle{turtle}
	w.Start()
	for w.IsRunning() {
		w.Tick()
	}
	got := turtle.pos
	want := coords.Pos{1, 1, 0}
	if got != want {
		t.Fatalf("got %v but want %v", got, want)
	}
}

func TestMovementBlocked(t *testing.T) {
	w := NewWorld(5)
	w.grid[coords.Pos{0, 1, 0}] = blocks.GetBlock(blocks.Bedrock)
	turtle := NewTurtle(coords.Pos{0, 0, 0}, w, coords.North).(*turtle)
	turtle.SetInfiniteFuel()
	turtle.SetProgram(func(tt Turtle) {
		tt.Forward()
		tt.TurnRight()
		tt.Forward()
	})
	w.Turtles = []Turtle{turtle}
	w.Start()
	for w.IsRunning() {
		w.Tick()
	}
	got := turtle.pos
	want := coords.Pos{1, 0, 0}
	if got != want {
		t.Fatalf("got %v but want %v", got, want)
	}
}

// Scenario: start facing a chest filled with coal
// Whenever fuel is almost out, take a few more coal
// Simulate for n steps and verify the amount of coal left + fuel level of turtle
func TestLoopWithRefuelStation(t *testing.T) {
	w := NewWorld(5)
	chest := &blocks.Chest{
		BaseBlock: blocks.BaseBlock{
			Type: blocks.ChestType,
		},
		Inventory: *blocks.NewInventory(16),
	}
	chest.Set(0, blocks.Stack{Type: blocks.Coal, Count: 64})
	w.Write(coords.Pos{0, 1, 0}, chest)
	turtle := NewTurtle(coords.Pos{0, 0, 0}, w, coords.North).(*turtle)
	turtle.inventory.Set(0, blocks.Stack{Type: blocks.Coal, Count: 5})
	turtle.SetProgram(func(tt Turtle) {
		for n := 0; n < 20; n++ {
			tt.Suck(2)
			tt.Refuel()
			for {
				tt.Up()
				tt.Down()
				if tt.GetFuelLevel() < 2 {
					break
				}
			}
		}
	})
	w.Turtles = []Turtle{turtle}
	w.Start()
	for w.IsRunning() {
		w.Tick()
	}
	got, want := turtle.pos, coords.Pos{0, 0, 0}
	if got != want {
		t.Errorf("got %v but want %v", got, want)
	}
	gotCoal, wantCoal := chest.Get(0).Count, 24
	if gotCoal != wantCoal {
		t.Errorf("Coal in chest: got %v but want %v", gotCoal, wantCoal)
	}
	gotCoal, wantCoal = turtle.inventory.Get(0).Count, 0
	if gotCoal != wantCoal {
		t.Errorf("Coal in turtle inv: got %v but want %v", gotCoal, wantCoal)
	}
}
