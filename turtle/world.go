package turtle

import (
	"fmt"
	"sync"

	"github.com/deosjr/TurtleSimulator/blocks"
	"github.com/deosjr/TurtleSimulator/coords"
)

// lives here until i promote it to yet another package

// convention: world minimal coord is 0,0,0 and max dim,dim,dim
type World struct {
	Dim     int
	Turtles []Turtle

	mu   sync.Mutex
	grid map[coords.Pos]blocks.Block
}

func NewWorld(dim int) *World {
	return &World{
		Dim:  dim,
		grid: map[coords.Pos]blocks.Block{},
	}
}

func (w *World) Tick() {
	for _, t := range w.Turtles {
		if !t.IsRunning() {
			continue
		}
		t.Tick()
	}
	for _, t := range w.Turtles {
		if !t.IsRunning() {
			continue
		}
		t.Tack()
	}
}

func (w *World) Start() {
	for _, t := range w.Turtles {
		go t.Run()
		for !t.IsRunning() {
		}
	}
}

func (w *World) IsRunning() bool {
	for _, t := range w.Turtles {
		if t.IsRunning() {
			return true
		}
	}
	return false
}

func (w *World) NumBlocks() int {
	return len(w.grid)
}

// TODO: remove?
func (w *World) Grid() map[coords.Pos]blocks.Block {
	return w.grid
}

func (w *World) Read(p coords.Pos) (blocks.Block, bool) {
	w.mu.Lock()
	v, ok := w.grid[p]
	w.mu.Unlock()
	return v, ok
}

func (w *World) Write(p coords.Pos, b blocks.Block) {
	w.mu.Lock()
	w.grid[p] = b
	w.mu.Unlock()
}

func (w *World) Delete(p coords.Pos) {
	w.mu.Lock()
	delete(w.grid, p)
	w.mu.Unlock()
}

// combines read, write and delete within 1 mutex lock/unlock
func (w *World) Move(from, to coords.Pos) (blocks.Block, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	block, ok := w.grid[from]
	if !ok {
		return block, fmt.Errorf("no block found")
	}
	_, ok = w.grid[to]
	if ok {
		return block, fmt.Errorf("block in position")
	}
	delete(w.grid, from)
	w.grid[to] = block
	return block, nil
}
