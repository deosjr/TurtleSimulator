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
    turtleHeading := t.(*turtle).heading

    for t.IsRunning() {
	    // only render if world has changed since last tick
        // todo: change in turtle position vs adding/removing blocks should
        // result in different optimisations (rebuilding bvh or not, for example)
        turtleMoved := turtlePos != t.(*turtle).pos
        turtleRotated := turtleHeading != t.(*turtle).heading
        blockPlacedOrRemoved := len(w.grid) != numBlocks

        if !turtleMoved && !turtleRotated && !blockPlacedOrRemoved {
            v.VisualiseUnchanged(w)
		    w.tick <- true
		    fmt.Println(turtlePos, turtleHeading)
		    <-w.tack
            continue
        }

        numBlocks = len(w.grid)
        turtlePos = t.(*turtle).pos
        turtleHeading = t.(*turtle).heading

        v.Visualise(w)
	    // send tick update to turtle and await yield
	    // todo: abort if turtle takes too long
	    w.tick <- true
	    fmt.Println(turtlePos, turtleHeading)
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

// experimental: to be taken into grayt if working
type cubetexture struct {
	colorFunc   func(textureSpace m.Vector) m.Color
}

func (t cubetexture) GetColor(si *m.SurfaceInteraction) m.Color {
	uv := texfunc(si)
	return t.colorFunc(uv)
}

// copied from model/texture:62 TriangleMeshUVFunc
func texfunc(si *m.SurfaceInteraction) m.Vector {
    ulhc := m.Vector{0,0,0}
    urhc := m.Vector{1,0,0}
    llhc := m.Vector{0,1,0}
    lrhc := m.Vector{1,1,0}
    tr := si.GetObject().(m.TriangleInMesh)
    p := si.UntransformedPoint
    l0, l1, l2 := tr.Barycentric(p)
    p0, p1, p2 := tr.PointIndices()
    var uv0, uv1, uv2 m.Vector
    var z float32
    // mapping cube coords: each face to square 0,0 - 1,1
    // alternative is to accurately map onto a cross cube layout
    switch {
    case p0 == 0 &&  p1 == 1 && p2 == 3:
        uv0, uv1, uv2 = ulhc, urhc, llhc
        z = 0
    case p0 == 1 && p1 == 2 && p2 == 3:
        uv0, uv1, uv2 = urhc, lrhc, llhc
        z = 0
    case p0 == 1 &&  p1 == 0 && p2 == 5:
        uv0, uv1, uv2 = ulhc, urhc, llhc
        z = 1
    case p0 == 0 && p1 == 4 && p2 == 5:
        uv0, uv1, uv2 = urhc, lrhc, llhc
        z = 1
    case p0 == 2 &&  p1 == 1 && p2 == 6:
        uv0, uv1, uv2 = ulhc, urhc, llhc
        z = 2
    case p0 == 1 && p1 == 5 && p2 == 6:
        uv0, uv1, uv2 = urhc, lrhc, llhc
        z = 2
    case p0 == 3 &&  p1 == 2 && p2 == 7:
        uv0, uv1, uv2 = ulhc, urhc, llhc
        z = 3
    case p0 == 2 && p1 == 6 && p2 == 7:
        uv0, uv1, uv2 = urhc, lrhc, llhc
        z = 3
    case p0 == 0 &&  p1 == 3 && p2 == 4:
        uv0, uv1, uv2 = ulhc, urhc, llhc
        z = 4
    case p0 == 3 && p1 == 7 && p2 == 4:
        uv0, uv1, uv2 = urhc, lrhc, llhc
        z = 4
    case p0 == 5 &&  p1 == 4 && p2 == 6:
        uv0, uv1, uv2 = ulhc, urhc, llhc
        z = 5
    case p0 == 4 && p1 == 7 && p2 == 6:
        uv0, uv1, uv2 = urhc, lrhc, llhc
        z = 5
    }
    uv := uv0.Times(l0).Add(uv1.Times(l1)).Add(uv2.Times(l2))
    // hide facenum in uv.Z coord
    uv.Z = z
    return uv
}

// instead of abusing z coord, this could be done with color values
// on each vertex in the mesh, then blending just like gettin uv values
func cf(uv m.Vector) m.Color {
    switch uv.Z {
    // front of the turtle is yellow
    case 1:
        return m.NewColor(250, 230, 60)
    default:
        return m.NewColor(255,0,0)
    }
}

func getturtleobj(cube m.AABB) m.Object {
    // copied from model/cube.go:128 cuboid tesselate
    f := func(p1,p2,p3,p4 int64) (m.Face, m.Face) {
        return m.Face{p1, p2, p4}, m.Face{p2, p3, p4}
    }
    pmin := cube.Pmin
	pmax := cube.Pmax

	p1 := m.Vector{pmin.X, pmax.Y, pmax.Z}
	p2 := m.Vector{pmax.X, pmax.Y, pmax.Z}
	p3 := m.Vector{pmax.X, pmax.Y, pmin.Z}
	p4 := m.Vector{pmin.X, pmax.Y, pmin.Z}

	p5 := m.Vector{pmin.X, pmin.Y, pmax.Z}
	p6 := m.Vector{pmax.X, pmin.Y, pmax.Z}
	p7 := m.Vector{pmax.X, pmin.Y, pmin.Z}
	p8 := m.Vector{pmin.X, pmin.Y, pmin.Z}
    vertices := []m.Vector{p1, p2, p3, p4, p5, p6, p7, p8}

	faces := make([]m.Face, 12)
	faces[0], faces[1] = f(0, 1, 2, 3)
	faces[2], faces[3] = f(1, 0, 4, 5)
	faces[4], faces[5] = f(2, 1, 5, 6)
	faces[6], faces[7] = f(3, 2, 6, 7)
	faces[8], faces[9] = f(0, 3, 7, 4)
	faces[10], faces[11] = f(5, 4, 7, 6)
    turtlemat := m.NewDiffuseMaterial(cubetexture{colorFunc: cf})
    turtleobj := m.NewTriangleMesh(vertices, faces, turtlemat)
    return turtleobj
}
// end copy

func (r *raytracer) Visualise(w *world) {

	r.scene.Objects = []m.Object{}

    // centering around the origin allows for easy rotations
	cube := m.NewAABB(m.Vector{-0.5, -0.5, -0.5}, m.Vector{0.5, 0.5, 0.5})

	for k, v := range w.grid {
		// z is up in turtle world, y is up in raytracing world
		// also left-right seem to be reversed _again_ (i thought id fixed that)
        // 0.5 is added to map to 0,0,0 through 1,1,1
		transform := m.Translate(m.Vector{float32(-k.x) + 0.5, float32(k.z) + 0.5, float32(k.y) + 0.5})
		var mat m.Material
		switch v.(type) {
		case Turtle:
			//mat = m.NewDiffuseMaterial(m.NewConstantTexture(m.NewColor(255, 0, 0)))
            switch v.(*turtle).heading {
            //case pos{0, 1, 0}: dont rotate when facing north
            case pos{1, 0, 0}:
                transform = transform.Mul(m.RotateY(math.Pi/2.0))
            case pos{0, -1, 0}:
                transform = transform.Mul(m.RotateY(math.Pi))
            case pos{-1, 0, 0}:
                transform = transform.Mul(m.RotateY(-math.Pi/2.0))
            }
		    shared := m.NewSharedObject(getturtleobj(cube), transform)
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

func (r *raytracer) VisualiseUnchanged(w *world) {
    render.AddToAVI(r.avi, r.film)
}

func (r *raytracer) Finalise() {
	render.SaveAVI(r.avi)
}
