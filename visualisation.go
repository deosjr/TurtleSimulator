package main

import (
	"fmt"
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
    "github.com/icza/mjpeg"
)

type visualiser interface {
    Visualise(*world)
    VisualiseUnchanged(*world)
    Finalise()
}

func Visualise(v visualiser, w *world) {
    // for now, exactly one turtle in the world
    t := w.turtles[0]
    go t.Run()

    v.Visualise(w)
    numBlocks := len(w.grid)
    turtlePos := t.(*turtle).pos
    for t.IsRunning() {
	    // only render if world has changed since last tick
        // todo: change in turtle position vs adding/removing blocks should
        // result in different optimisations (rebuilding bvh or not, for example)
        if turtlePos == t.(*turtle).pos && len(w.grid) == numBlocks {
            v.VisualiseUnchanged(w)
		    w.tick <- true
		    fmt.Println(t.(*turtle).pos)
		    <-w.tack
            continue
        }
        numBlocks = len(w.grid)
        turtlePos = t.(*turtle).pos
        v.Visualise(w)
	    // send tick update to turtle and await yield
	    // todo: abort if turtle takes too long
	    w.tick <- true
	    fmt.Println(t.(*turtle).pos)
	    <-w.tack
    }
    v.Visualise(w)
    v.Finalise()
}

type ascii struct {}

// prints grid from 0,0,0 at bottom left to dim,dim,dim at top right
// for each x,y coord, prints only the highest block, if any (top down view)
func printworld(w *world) {
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

func (ascii) Visualise(w *world) { printworld(w) }
func (ascii) VisualiseUnchanged(w *world) { printworld(w) }
func (ascii) Finalise() {}

type raytracer struct {
    // todo: hide this import in grayt
    avi  mjpeg.AviWriter
    film render.Film
    camera m.Camera
    scene  *m.Scene
}

func NewRaytracer() *raytracer {
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

    return &raytracer{
        avi:    avi,
        camera: camera,
        scene:  scene,
    }
}

func (r *raytracer) Visualise(w *world) {

	r.scene.Objects = []m.Object{}

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
		r.scene.Add(shared)
	}
	r.scene.Precompute()

	from, to := m.Vector{0, 2, -5}, m.Vector{0, 0, 10}
	r.camera.LookAt(from, to, m.Vector{0, 1, 0})

	params := render.Params{
		Scene:        r.scene,
		NumWorkers:   10,
		NumSamples:   10,
		AntiAliasing: true,
		TracerType:   m.WhittedStyle,
		//TracerType: m.PathNextEventEstimate,
	}
	r.film = render.Render(params)
	render.AddToAVI(r.avi, r.film)
}

func (r *raytracer) VisualiseUnchanged(w *world) {
    render.AddToAVI(r.avi, r.film)
}

func (r *raytracer) Finalise() {
	render.SaveAVI(r.avi)
}
