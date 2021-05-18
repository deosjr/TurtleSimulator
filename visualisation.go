package main

import (
	"fmt"
    "image"
    "image/color"
	"math"

    "github.com/deosjr/TurtleSimulator/blocks"
    "github.com/deosjr/TurtleSimulator/coords"
    "github.com/deosjr/TurtleSimulator/turtle"

	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GRayT/src/render"
	"github.com/deosjr/GenGeo/gen"
    "github.com/icza/mjpeg"
)

type visualiser interface {
    Visualise(*turtle.World)
    VisualiseMove(w *turtle.World, from, to coords.Pos)
    VisualiseUnchanged(*turtle.World)
    Finalise(entireRun bool)
}

// showEntireRun=false means we only want to show the end state
func visualise(v visualiser, w *turtle.World, showEntireRun bool) {
    // for now, we focus on one turtle only (this matters e.g. for visualising movement frames)
    t := w.Turtles[0]
    w.Start()

    if showEntireRun {
        v.Visualise(w)
    }
    numBlocks := w.NumBlocks()
    turtlePos := t.GetPos()
    turtleHeading := t.GetHeading()

    for w.IsRunning() {
	    // only render if world has changed since last tick
        // todo: change in turtle position vs adding/removing blocks should
        // result in different optimisations (rebuilding bvh or not, for example)
        turtleMoved := turtlePos != t.GetPos()
        turtleRotated := turtleHeading != t.GetHeading()
        blockPlacedOrRemoved := w.NumBlocks() != numBlocks

        if !turtleMoved && !turtleRotated && !blockPlacedOrRemoved {
            // when turtle detects or otherwise yields without changing
            if showEntireRun {
                v.VisualiseUnchanged(w)
            }
		    w.Tick()
		    fmt.Println(turtlePos, coords.HeadingString(turtleHeading))
            continue
        }

        if showEntireRun {
            if turtleMoved {
                v.VisualiseMove(w, turtlePos, t.GetPos())
            } else {
                v.Visualise(w)
            }
        }

        numBlocks = w.NumBlocks()
        turtlePos = t.GetPos()
        turtleHeading = t.GetHeading()

	    // send tick update to turtles and await yield
	    // todo: abort if turtle takes too long
	    w.Tick()
		fmt.Println(turtlePos, coords.HeadingString(turtleHeading))
    }
    if showEntireRun && turtlePos != t.GetPos() {
        v.VisualiseMove(w, turtlePos, t.GetPos())
    }
    v.Visualise(w)
    v.Finalise(showEntireRun)
}

func Visualise(v visualiser, w *turtle.World) {
    visualise(v, w, true)
}

func VisualiseEndState(v visualiser, w *turtle.World) {
    visualise(v, w, false)
}

type ascii struct {}

// prints grid from 0,0,0 at bottom left to dim,dim,dim at top right
// for each x,y coord, prints only the highest block, if any (top down view)
func printworld(w *turtle.World) {
	for y := w.Dim; y >= 0; y-- {
	Loop:
		for x := 0; x < 5; x++ {
			for z := w.Dim; z >= 0; z-- {
				b, ok := w.Read(coords.Pos{x, y, 0})
				if !ok {
					continue
				}
				switch t := b.(type) {
				case turtle.Turtle:
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

func (ascii) Visualise(w *turtle.World) { printworld(w) }
func (ascii) VisualiseMove(w *turtle.World, from, to coords.Pos) { printworld(w) }
func (ascii) VisualiseUnchanged(w *turtle.World) { printworld(w) }
func (ascii) Finalise(_ bool) {}

type raytracer struct {
    // todo: hide this import in grayt
    avi  mjpeg.AviWriter
    film render.Film
    camera m.Camera
    scene  *m.Scene
    move bool
    followTurtle bool
}

func NewRaytracer(visualiseMove, follow bool) *raytracer {
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
        move:   visualiseMove,
        followTurtle: follow,
    }
}

func getturtleobj() m.Object {
    img := image.NewRGBA(image.Rect(0,0,4,3))
    front := color.RGBA{240, 200, 60, 255}
    cube := color.RGBA{255, 0, 0, 255}
    img.Set(0,1,front)
    img.Set(1,0,cube)
    img.Set(1,1,cube)
    img.Set(1,2,cube)
    img.Set(2,1,cube)
    img.Set(3,1,cube)
    return m.CubeMesh(0.7, img)
}

func (r *raytracer) Visualise(w *turtle.World) {
    r.visualise(w, 0, 0, 0)
}

func (r *raytracer) visualise(w *turtle.World, dx, dy, dz float32) {

	r.scene.Objects = []m.Object{}

    // centering around the origin allows for easy rotations
	cube := m.NewAABB(m.Vector{-0.5, -0.5, -0.5}, m.Vector{0.5, 0.5, 0.5})
    stairpoints := []m.Vector{
        {-0.5, -0.5, -0.5},
        {0.5, -0.5, -0.5},
        {0.5, 0, -0.5},
        {0, 0, -0.5},
        {0, 0.5, -0.5},
        {-0.5, 0.5, -0.5},
    }

    var turtlepos m.Vector

	for k, v := range w.Grid() {
		// z is up in turtle world, y is up in raytracing world
        // 0.5 is added to map to 0,0,0 through 1,1,1
		transform := m.Translate(m.Vector{float32(-k.X) + 0.5, float32(k.Z) + 0.5, float32(k.Y) + 0.5})
		var mat m.Material
		switch v.GetType() {
		case blocks.Turtle:
		    transform = m.Translate(m.Vector{float32(-k.X) + 0.5 + dx, float32(k.Z) + 0.5 + dz, float32(k.Y) + 0.5 + dy})
            switch v.GetHeading() {
            //case pos{0, 1, 0}: dont rotate when facing north
            case coords.East:
                transform = transform.Mul(m.RotateY(math.Pi/2.0))
            case coords.South:
                transform = transform.Mul(m.RotateY(math.Pi))
            case coords.West:
                transform = transform.Mul(m.RotateY(-math.Pi/2.0))
            }
		    shared := m.NewSharedObject(getturtleobj(), transform)
		    r.scene.Add(shared)
            turtlepos = m.Vector{float32(-k.X) + 0.5 + dx, float32(k.Z) + 0.5 + dz, float32(k.Y) + 0.5 + dy}
            continue
		case blocks.Grass:
			mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(0, 255, 0)))
		case blocks.Stone:
			mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(150, 150, 150)))
		case blocks.Log:
			mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(90, 60, 10)))
		case blocks.Planks:
			mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(150, 100, 20)))
		case blocks.Brick:
			mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(190, 190, 190)))
		case blocks.CobbleSlab, blocks.BrickSlab: // TODO switch colors between them
		    transform = m.Translate(m.Vector{float32(-k.X) + 0.5, float32(k.Z) + 0.25, float32(k.Y) + 0.5})
            if v.(blocks.BaseBlock).Flipped {
		        transform = m.Translate(m.Vector{float32(-k.X) + 0.5, float32(k.Z) + 0.75, float32(k.Y) + 0.5})
            }
			mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(190, 190, 190)))
	        slab := m.NewAABB(m.Vector{-0.5, -0.25, -0.5}, m.Vector{0.5, 0.25, 0.5})
		    block := m.NewCuboid(slab, mat).Tesselate()
		    shared := m.NewSharedObject(m.NewTriangleComplexObject(block), transform)
		    r.scene.Add(shared)
            continue
        case blocks.Stairs:
			mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(150, 150, 150)))
            stairsobj := gen.ExtrudeSolidFace(stairpoints, m.Vector{0,0,1}, mat)
            switch v.GetHeading() {
            //case pos{0, 1, 0}: dont rotate when facing north
            case coords.East:
                transform = transform.Mul(m.RotateY(math.Pi/2.0))
            case coords.South:
                transform = transform.Mul(m.RotateY(math.Pi))
            case coords.West:
                transform = transform.Mul(m.RotateY(-math.Pi/2.0))
            }
            if v.(blocks.BaseBlock).Flipped {
                transform = transform.Mul(m.RotateX(math.Pi))
            }
            shared := m.NewSharedObject(stairsobj, transform)
            r.scene.Add(shared)
            continue
		default:
			mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(0, 0, 255)))
		}
		block := m.NewCuboid(cube, mat).Tesselate()
		shared := m.NewSharedObject(m.NewTriangleComplexObject(block), transform)
		r.scene.Add(shared)
	}
	r.scene.Precompute()

    if !r.followTurtle {
	    from, to := m.Vector{10, 15, -5}, m.Vector{0, 0, 10}
	    //from, to := m.Vector{5, 8, -5}, m.Vector{5, 8, 10}
	    r.camera.LookAt(from, to, m.Vector{0, 1, 0})
    } else {
	    to := turtlepos
        from := m.Vector{to.X, to.Y + 2, to.Z - 5}
	    r.camera.LookAt(from, to, m.Vector{0, 1, 0})
    }

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
    r.film.SaveAsPNG("turtle.png")
}

func (r *raytracer) VisualiseMove(w *turtle.World, from, to coords.Pos) {
    if !r.move {
        return
    }
    dif := to.Sub(from)
    // i will regret flipping the x-axis. here's the first regret
    stepx := -float32(dif.X) / 8.0
    stepy := float32(dif.Y) / 8.0
    stepz := float32(dif.Z) / 8.0
    for i:=8; i>0; i-- {
        fi := float32(-i)
        r.visualise(w, fi*stepx, fi*stepy, fi*stepz)
    }
}

func (r *raytracer) VisualiseUnchanged(w *turtle.World) {
    render.AddToAVI(r.avi, r.film)
}

func (r *raytracer) Finalise(entireRun bool) {
    if entireRun {
	    render.SaveAVI(r.avi)
    }
    r.avi.Close()
}
