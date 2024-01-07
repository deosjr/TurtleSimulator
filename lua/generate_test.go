package main

import (
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func TestFirstPassSimple(t *testing.T) {
	src := `package main
    
	import "github.com/deosjr/TurtleSimulator/turtle"

    // comment
    func Forward(t turtle.Turtle) {
        t.forward()
        t.turnLeft()
        t.back()
    }`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	progs := programs(f)
	if len(progs) != 1 {
		t.Fatalf("Expected 1 program but found %d", len(progs))
	}
	got := generateFirstPass(progs[0])
	funcName := "forward"
	want := codeBlock{
		funcName: &funcName,
		lines: []line{
			{s: "turtle.forward()"},
			{s: "turtle.turnLeft()"},
			{s: "turtle.back()"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestFirstPassIfElse(t *testing.T) {
	src := `package main
    
	import "github.com/deosjr/TurtleSimulator/turtle"

    // comment
    func Forward(t turtle.Turtle) {
        if t.detect() {
            t.forward()
            t.turnLeft()
        } else {
            t.back()
        }
    }`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	progs := programs(f)
	if len(progs) != 1 {
		t.Fatalf("Expected 1 program but found %d", len(progs))
	}
	got := generateFirstPass(progs[0])
	funcName := "forward"
	want := codeBlock{
		funcName: &funcName,
		lines: []line{
			{s: "function () i = mem.condJump(i, 4, not turtle.detect()) end"},
			{s: "turtle.forward()"},
			{s: "turtle.turnLeft()"},
			{s: "function () i = mem.goto(i, 2) end"},
			{s: "turtle.back()"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestFirstPassFor(t *testing.T) {
	src := `package main
    
	import "github.com/deosjr/TurtleSimulator/turtle"

    // comment
    func Forward(t turtle.Turtle) {
        for n := 0; n < 2; n++ {
            t.Forward()
        }
    }`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	progs := programs(f)
	if len(progs) != 1 {
		t.Fatalf("Expected 1 program but found %d", len(progs))
	}
    gensymFunc = func() string { return "n" }
	got := generateFirstPass(progs[0])
	funcName := "forward"
	want := codeBlock{
		funcName: &funcName,
		lines: []line{
            {s: "function () if state.n == nil then state.n = 0 end i = mem.condJump(i, 3, state.n >= 2) end"},
			{s: "function () turtle.forward(); state.n = state.n+1"},
            {s: "function () i = mem.goto(i, -2) end"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestFirstPassWhileTrue(t *testing.T) {
	src := `package main
    
	import "github.com/deosjr/TurtleSimulator/turtle"

    // comment
    func Forward(t turtle.Turtle) {
        for {
            t.Forward()
        }
    }`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	progs := programs(f)
	if len(progs) != 1 {
		t.Fatalf("Expected 1 program but found %d", len(progs))
	}
	got := generateFirstPass(progs[0])
	funcName := "forward"
	want := codeBlock{
		funcName: &funcName,
		lines: []line{
			{s: "turtle.forward()"},
			{s: "function () i = mem.goto(i, -2) end"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestFirstPassWhileCond(t *testing.T) {
	src := `package main
    
	import "github.com/deosjr/TurtleSimulator/turtle"

    // comment
    func Forward(t turtle.Turtle) {
        for !t.Detect() {
            t.Forward()
        }
    }`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	progs := programs(f)
	if len(progs) != 1 {
		t.Fatalf("Expected 1 program but found %d", len(progs))
	}
	got := generateFirstPass(progs[0])
	funcName := "forward"
	want := codeBlock{
		funcName: &funcName,
		lines: []line{
			{s: "function () i = mem.condJump(i, 3, turtle.detect()) end"},
			{s: "turtle.forward()"},
			{s: "function () i = mem.goto(i, -2) end"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestFirstPassNestedFunc(t *testing.T) {
	src := `package main
    
	import "github.com/deosjr/TurtleSimulator/turtle"

    // comment
    func Forward(t turtle.Turtle) {
        t.Forward()
        sidestepRight(t)
        t.Back()
    }

    func sidestepRight(t turtle.Turtle) {
        t.TurnRight()
        t.Forward()
        t.TurnLeft()
    }`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	progs := programs(f)
	if len(progs) != 1 {
		t.Fatalf("Expected 1 program but found %d", len(progs))
	}
	got := generateFirstPass(progs[0])
	funcName := "forward"
	want := codeBlock{
		funcName: &funcName,
		lines: []line{
			{s: "turtle.forward()"},
			// TODO: instead of inlining the function, use a callstack and goto
			{s: "turtle.turnRight()"},
			{s: "turtle.forward()"},
			{s: "turtle.turnLeft()"},
			{s: "turtle.back()"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}

func TestGenerateSimpleProgram(t *testing.T) {
	src := `package main
    
	import "github.com/deosjr/TurtleSimulator/turtle"

    // comment
    func Forward(t turtle.Turtle) {
        t.forward()
    }`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	progs := programs(f)
	if len(progs) != 1 {
		t.Fatalf("Expected 1 program but found %d", len(progs))
	}
	got := generateSimpleProgram(progs[0])
	want := `-- comment
function forward()
    turtle.forward()
end`
	if got != want {
		t.Fatalf("Expected %s but got %s", want, got)
	}
}

func TestGenerateProgram(t *testing.T) {
	src := `package main
    
	import "github.com/deosjr/TurtleSimulator/turtle"

    // comment
    func Forward(t turtle.Turtle) {
        t.forward()
    }`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	progs := programs(f)
	if len(progs) != 1 {
		t.Fatalf("Expected 1 program but found %d", len(progs))
	}
	got := generateProgram(progs[0])
	want := `-- comment
local state = mem.startFromMemory("forward")
if not state then
    state = {i=0}
end
local i = state.i

local stop = false -- used to communicate key Q pressed

action = {
    turtle.forward(),
}

function main()
    while not stop do
        -- table is 1-based, modulo arithmetic is 0-based
        action[i+1]()
        i = ((i + 1) % #action)
        state.i = i
        mem.writeMemory("forward", state)
    end
end

function keyInterrupt()
    while true do
        local event, key, isHeld = os.pullEvent("key")
        if key == keys.q then
            stop = true
            break
        end
    end
end

parallel.waitForAny(main, keyInterrupt)`
	if got != want {
		t.Fatalf("Expected %s but got %s", want, got)
	}
}
