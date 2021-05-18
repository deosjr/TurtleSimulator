package main

import (
    "github.com/deosjr/TurtleSimulator/scene"
)


// TODO: a game tick in minecraft is 1/20 second
// from dan200 computercraft java code:
// each animation takes 8 ticks to complete unless otherwise specified.

func main() {
    w := scene.Ziggurat()

    vis := NewRaytracer(false, false)
    //vis := ascii{}
    //Visualise(vis, w)
    VisualiseEndState(vis, w)
}
