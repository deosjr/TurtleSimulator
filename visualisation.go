package main

import (
	"fmt"
    "image"
    "image/color"
	"math"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
    "github.com/icza/mjpeg"
)

type visualiser interface {
    Visualise(*world)
    VisualiseMove(w *world, from, to pos)
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
    turtleHeading := t.(*turtle).heading

    for t.IsRunning() {
	    // only render if world has changed since last tick
        // todo: change in turtle position vs adding/removing blocks should
        // result in different optimisations (rebuilding bvh or not, for example)
        turtleMoved := turtlePos != t.(*turtle).pos
        turtleRotated := turtleHeading != t.(*turtle).heading
        blockPlacedOrRemoved := len(w.grid) != numBlocks

        if !turtleMoved && !turtleRotated && !blockPlacedOrRemoved {
            // when turtle detects or otherwise yields without changing
            v.VisualiseUnchanged(w)
		    w.tick <- true
		    fmt.Println(turtlePos, turtleHeading)
		    <-w.tack
            continue
        }


        if turtleMoved {
            v.VisualiseMove(w, turtlePos, t.(*turtle).pos)
        } else {
            v.Visualise(w)
        }

        numBlocks = len(w.grid)
        turtlePos = t.(*turtle).pos
        turtleHeading = t.(*turtle).heading

	    // send tick update to turtle and await yield
	    // todo: abort if turtle takes too long
	    w.tick <- true
	    fmt.Println(turtlePos, turtleHeading)
	    <-w.tack
    }
    if turtlePos != t.(*turtle).pos {
        v.VisualiseMove(w, turtlePos, t.(*turtle).pos)
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
func (ascii) VisualiseMove(w *world, from, to pos) { printworld(w) }
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

func getturtleobj() m.Object {
    /*
    // img is mapped onto the cube as follows:
        |---| 
        | 2 |
    |---|---|---|---|
    | 1 | 3 | 5 | 6 |
    |---|---|---|---|
        | 4 |
        |---|
    // 1 is the front of the cube, facing in +Z direction
    // 2 is right, 3 is bottom, 4 is left, 5 is back, 6 is top
    // each square has topleft at topleft, so oriented the same
    // therefore NOT neatly wrapping the whole cross around the cube!
    // works for arbitrary resolution as long as aspect ratio is 4:3
    */
    img := image.NewRGBA(image.Rect(0,0,4,3))
    front := color.RGBA{240, 200, 60, 255}
    cube := color.RGBA{255, 0, 0, 255}
    img.Set(0,1,front)
    img.Set(1,0,cube)
    img.Set(1,1,cube)
    img.Set(1,2,cube)
    img.Set(2,1,cube)
    img.Set(3,1,cube)

    // copied from model/cube.go:128 cuboid tesselate
    f := func(p0,p1,p2,p3 int64) (m.Face, m.Face) {
        return m.Face{p0, p2, p1}, m.Face{p1, p2, p3}
    }
    // unit cube centered around origin
    var min, max float32 = -0.5, 0.5

    // see ilkinulas.github.io/development/unity/2016/05/06/uv-mapping.html for details
	p0 := m.Vector{max, max, max}
	p1 := m.Vector{max, min, max}
	p2 := m.Vector{min, max, max}
	p3 := m.Vector{min, min, max}
	p4 := m.Vector{max, min, min}
	p5 := m.Vector{min, min, min}
	p6 := m.Vector{max, max, min}
	p7 := m.Vector{min, max, min}
    vertices := []m.Vector{p0, p1, p2, p3, p4, p5, p6, p7, p0, p2, p0, p6, p2, p7}

	faces := make([]m.Face, 12)
	faces[0], faces[1] = f(0, 1, 2, 3)
	faces[2], faces[3] = f(10, 11, 1, 4)
	faces[4], faces[5] = f(1, 4, 3, 5)
	faces[6], faces[7] = f(3, 5, 12, 13)
	faces[8], faces[9] = f(4, 6, 5, 7)
	faces[10], faces[11] = f(6, 8, 7, 9)
    uvmap := map[int64]m.Vector{
        0:  m.Vector{0, 2.0/3.0, 0},
        1:  m.Vector{0.25, 2.0/3.0, 0},
        2:  m.Vector{0, 1.0/3.0, 0},
        3:  m.Vector{0.25, 1.0/3.0, 0},
        4:  m.Vector{0.5, 2.0/3.0, 0},
        5:  m.Vector{0.5, 1.0/3.0, 0},
        6:  m.Vector{0.75, 2.0/3.0, 0},
        7:  m.Vector{0.75, 1.0/3.0, 0},
        8:  m.Vector{1, 2.0/3.0, 0},
        9:  m.Vector{1, 1.0/3.0, 0},
        10: m.Vector{0.25, 1, 0},
        11: m.Vector{0.5, 1, 0},
        12: m.Vector{0.25, 0, 0},
        13: m.Vector{0.5, 0, 0},
    }

    turtlemat := m.NewDiffuseMaterial(m.NewImageTexture(img, m.TriangleMeshUVFunc))
    turtleobj := m.NewTriangleMesh(vertices, faces, turtlemat)
    turtleobj.(*m.TriangleMesh).UV = uvmap
    return turtleobj
}

func (r *raytracer) Visualise(w *world) {
    r.visualise(w, 0, 0, 0)
}

func (r *raytracer) visualise(w *world, dx, dy, dz float32) {

	r.scene.Objects = []m.Object{}

    // centering around the origin allows for easy rotations
	cube := m.NewAABB(m.Vector{-0.5, -0.5, -0.5}, m.Vector{0.5, 0.5, 0.5})

	for k, v := range w.grid {
		// z is up in turtle world, y is up in raytracing world
        // 0.5 is added to map to 0,0,0 through 1,1,1
		transform := m.Translate(m.Vector{float32(-k.x) + 0.5, float32(k.z) + 0.5, float32(k.y) + 0.5})
		var mat m.Material
		switch v.(type) {
		case Turtle:
			//mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(255, 0, 0)))
		    transform = m.Translate(m.Vector{float32(-k.x) + 0.5 + dx, float32(k.z) + 0.5 + dz, float32(k.y) + 0.5 + dy})
            switch v.(*turtle).heading {
            //case pos{0, 1, 0}: dont rotate when facing north
            case pos{1, 0, 0}:
                transform = transform.Mul(m.RotateY(math.Pi/2.0))
            case pos{0, -1, 0}:
                transform = transform.Mul(m.RotateY(math.Pi))
            case pos{-1, 0, 0}:
                transform = transform.Mul(m.RotateY(-math.Pi/2.0))
            }
		    shared := m.NewSharedObject(getturtleobj(), transform)
		    r.scene.Add(shared)
            continue
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

func (r *raytracer) VisualiseMove(w *world, from, to pos) {
    dif := to.sub(from)
    // i will regret flipping the x-axis. here's the first regret
    stepx := -float32(dif.x) / 8.0
    stepy := float32(dif.y) / 8.0
    stepz := float32(dif.z) / 8.0
    for i:=8; i>0; i-- {
        fi := float32(-i)
        r.visualise(w, fi*stepx, fi*stepy, fi*stepz)
    }
}

func (r *raytracer) VisualiseUnchanged(w *world) {
    render.AddToAVI(r.avi, r.film)
}

func (r *raytracer) Finalise() {
	render.SaveAVI(r.avi)
}
