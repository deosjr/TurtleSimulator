package main

import (
	"fmt"
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
)

type pos struct {
	x, y, z int
}

func (p pos) add(q pos) pos {
	return pos{p.x + q.x, p.y + q.y, p.z + q.z}
}

func (p pos) sub(q pos) pos {
	return pos{p.x - q.x, p.y - q.y, p.z - q.z}
}

// convention: world minimal coord is 0,0,0 and max dim,dim,dim
type world struct {
	grid map[pos]Block
	dim  int
}

// prints grid from 0,0,0 at bottom left to dim,dim,dim at top right
// for each x,y coord, prints only the highest block, if any (top down view)
func (w *world) print() {
	for y := w.dim; y >= 0; y-- {
	Loop:
		for x := 0; x < 5; x++ {
			for z := w.dim; z >= 0; z-- {
				b, ok := w.grid[pos{x, y, 0}]
				if !ok {
					continue
				}
				switch t := b.(type) {
				case Turtle:
					fmt.Print(t.String())
				default:
					fmt.Print("x")
				}
				continue Loop
			}
			fmt.Print(".")
		}
		fmt.Println()
	}
}

type Block interface{}
type block struct{}
type grass struct{}
type stone struct{}

// TODO: a game tick in minecraft is 1/20 second
// from dan200 computercraft java code:
// each animation takes 8 ticks to complete unless otherwise specified.

func main() {
	w := &world{grid: map[pos]Block{}, dim: 5}
	//w.grid[pos{0, 3, 0}] = block{}
	//w.grid[pos{2, 2, 0}] = block{}
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			w.grid[pos{x, y, -1}] = grass{}
		}
	}
	tick := make(chan bool, 1)
	tack := make(chan bool, 1)
	t := NewTurtle(pos{0, 0, 0}, w, tick, tack)
	// note: right now inner t shadows outer t...
	/*
		t.SetProgram(func(t Turtle) {
			for !t.Detect() {
				t.Forward()
			}
			t.TurnRight()
			for !t.Detect() {
				t.Forward()
			}
		})
	*/
	// wall building program
	t.SetProgram(func(t Turtle) {
		sidestepRight := func() {
			t.TurnRight()
			t.Forward()
			t.TurnLeft()
		}
		sidestepLeft := func() {
			t.TurnLeft()
			t.Forward()
			t.TurnRight()
		}
		if t.Detect() {
			sidestepRight()
			t.Forward()
			t.TurnLeft()
		} else {
			t.Forward()
			t.TurnLeft()
			t.Back()
		}
		for i := 0; i < 2; i++ {
			t.Place()
			t.TurnRight()
			t.Place()
			for j := 0; j < 2; j++ {
				sidestepRight()
				t.Place()
			}
			t.Up()
			t.Place()
			for j := 0; j < 2; j++ {
				sidestepLeft()
				t.Place()
			}
			t.TurnLeft()
			t.Place()
			t.Up()
		}
		t.Place()
		t.Up()
		t.Place()
		t.Down()
		t.TurnRight()
		t.Place()
		t.PlaceUp()

		for i := 0; i < 2; i++ {
			sidestepRight()
			t.Place()
			t.PlaceUp()
		}
		t.Down()
		t.TurnRight()
		t.PlaceUp()
		t.Back()
		t.Back()
		t.TurnRight()
		t.TurnRight()
		t.PlaceUp()
		for i := 0; i < 3; i++ {
			t.Down()
			t.Back()
		}
		t.TurnRight()
		t.Back()
	})
	go t.Run()

	// use grayt to render avi of turtle moving
	m.SIMD_ENABLED = true
	var width, height uint = 800, 600
	camera := m.NewPerspectiveCamera(width, height, 0.5*math.Pi)
	avi := render.NewAVI("turtle.avi", 800, 600)

	scene := m.NewScene(camera)
	/*
		radmat := m.NewRadiantMaterial(m.NewConstantTexture(m.NewColor(176, 237, 255)))
		skybox := m.NewCuboid(m.NewAABB(m.Vector{-1000, -1000, -1000}, m.Vector{1000, 1000, 1000}), radmat)
		triangles := skybox.TesselateInsideOut()
		skyboxObject := m.NewTriangleComplexObject(triangles)
		scene.Add(skyboxObject)
		scene.Emitters = triangles
	*/
	//pointLight := m.NewPointLight(m.Vector{0, 10, -100}, m.NewColor(255, 255, 255), 50000000)
	//scene.AddLights(pointLight)
	l1 := m.NewDistantLight(m.Vector{-1, -1, 1}, m.NewColor(255, 255, 255), 20)
	scene.AddLights(l1)
	m.SetBackgroundColor(m.NewColor(15, 200, 215))

	// todo: skips very last action atm
	// todo: only render if world has changed since last tick
	for t.IsRunning() {
		scene.Objects = []m.Object{}

		cube := m.NewAABB(m.Vector{0, 0, 0}, m.Vector{1, 1, 1})
		for k, v := range w.grid {
			// z is up in turtle world, y is up in raytracing world
			// also left-right seem to be reversed _again_ (i thought id fixed that)
			transform := m.Translate(m.Vector{float32(-k.x), float32(k.z), float32(k.y)})
			var mat m.Material
			switch v.(type) {
			case Turtle:
				mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(255, 0, 0)))
			case grass:
				mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(0, 255, 0)))
			case stone:
				mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(150, 150, 150)))
			default:
				mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(0, 0, 255)))
			}
			block := m.NewCuboid(cube, mat).Tesselate()
			shared := m.NewSharedObject(m.NewTriangleComplexObject(block), transform)
			scene.Add(shared)
		}
		scene.Precompute()

		from, to := m.Vector{0, 2, -5}, m.Vector{0, 0, 10}
		camera.LookAt(from, to, m.Vector{0, 1, 0})

		params := render.Params{
			Scene:        scene,
			NumWorkers:   10,
			NumSamples:   10,
			AntiAliasing: true,
			TracerType:   m.WhittedStyle,
			//TracerType: m.PathNextEventEstimate,
		}
		film := render.Render(params)
		render.AddToAVI(avi, film)
		// send tick update to turtle and await yield
		// todo: abort if turtle takes too long
		tick <- true
		fmt.Println(t.(*turtle).pos)
		<-tack
	}
	render.SaveAVI(avi)
}
