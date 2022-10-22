package lib

import (
	"github.com/deosjr/TurtleSimulator/blocks"
	"github.com/deosjr/TurtleSimulator/turtle"
)

// helper funcs, should later on also include selecting from inv if possible
func Place(t turtle.Turtle, bt blocks.Blocktype) bool {
	t.SetInventory(bt)
	return t.Place()
}
func PlaceUp(t turtle.Turtle, bt blocks.Blocktype) bool {
	t.SetInventory(bt)
	return t.PlaceUp()
}
func PlaceDown(t turtle.Turtle, bt blocks.Blocktype) bool {
	t.SetInventory(bt)
	return t.PlaceDown()
}
