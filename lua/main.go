package main

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// TODO: parse go turtle programs and output lua files :)

// type Program func(Turtle)
// convention: the Turtle is always called 't'

// first pass parsing is syntax -> codeblocks
// then we put all that together in lua 'assembler'
type codeBlock struct {
	funcName *string
	lines    []line
}

// TODO can include variables for goto etc
type line struct {
	s string
}

func main() {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	programspkg := pkgs[0]
	fmt.Println(programspkg.GoFiles)
	for i := 0; i < len(programspkg.GoFiles); i++ {
		//file := programspkg.GoFiles[i]
		syntax := programspkg.Syntax[i]
		for _, f := range programs(syntax) {
			// TODO: for now, only run for wallbuildfunc
			if f.Name.Name != "Wallbuildfunc" {
				continue
			}
			fmt.Println(generateProgram(f))
		}
	}
}

func programs(f *ast.File) []*ast.FuncDecl {
	programs := []*ast.FuncDecl{}
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// only consider exported funcs
		if !unicode.IsUpper(rune(fd.Name.Name[0])) {
			continue
		}
		if fd.Type.Results != nil {
			continue
		}
		if fd.Type.Params == nil || len(fd.Type.Params.List) != 1 {
			continue
		}
		paramType := fd.Type.Params.List[0].Type
		pt, ok := paramType.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		x := pt.X.(*ast.Ident).Name
		sel := pt.Sel.Name
		if x == "turtle" && sel == "Turtle" {
			// found a Program
			programs = append(programs, fd)
		}
	}
	return programs
}

type generator struct {
	b strings.Builder
}

func (g *generator) Writef(format string, args ...any) {
	fmt.Fprintf(&g.b, format, args...)
}

func generateSimpleProgram(fd *ast.FuncDecl) string {
	g := &generator{}
	/*
	   if fd.Doc != nil {
	       for _, line := range fd.Doc.List {
	           g.b.WriteString(line.Text)
	       }
	   }
	*/
	g.Writef("-- comment\n")
	name := strings.ToLower(fd.Name.Name)
	g.Writef("function %s()\n", name)
	for _, stmt := range fd.Body.List {
		ex, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		cx, ok := ex.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		sx, ok := cx.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		g.Writef("    %s\n", printLuaFunc(sx.X.(*ast.Ident).Name, sx.Sel.Name, cx.Args))
	}
	g.Writef("end")
	return g.b.String()
}

func generateProgram(fd *ast.FuncDecl) string {
	g := &generator{}
	name := strings.ToLower(fd.Name.Name)
	g.Writef(`-- comment
local state = mem.startFromMemory(%q)
if not state then
    state = {i=0}
end
local i = state.i

local stop = false -- used to communicate key Q pressed

action = {
`, name)
	for _, stmt := range fd.Body.List {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			cx, ok := s.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			switch f := cx.Fun.(type) {
			case *ast.SelectorExpr:
				g.Writef("    %s,\n", printLuaFunc(f.X.(*ast.Ident).Name, f.Sel.Name, cx.Args))
			case *ast.Ident:
				fmt.Printf("%v\n", f)
			}
		case *ast.ForStmt:
			fmt.Println("TODO FOR: ", s)
		case *ast.IfStmt:
			fmt.Println("TODO IFELSE: ", s)
		}
	}
	g.Writef(`}

function main()
    while not stop do
        -- table is 1-based, modulo arithmetic is 0-based
        action[i+1]()
        i = ((i + 1) %% #action)
        state.i = i
        mem.writeMemory(%q, state)
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

parallel.waitForAny(main, keyInterrupt)`, name)
	return g.b.String()
}

func printLuaFunc(a, b string, args []ast.Expr) string {
	b = strings.ToLower(string(b[0])) + b[1:]
	if a == "t" {
		return fmt.Sprintf("turtle." + b + "()")
	}
	if a == "lib" {
		return fmt.Sprintf("lib." + b + "()")
	}
	if a == "mem" {
		return fmt.Sprintf("mem." + b + "()")
	}
	return fmt.Sprintf("UNKNOWN<%s.%s>", a, b)
}

func generateFirstPass(fd *ast.FuncDecl) codeBlock {
	name := strings.ToLower(fd.Name.Name)
	return codeBlock{
		funcName: &name,
		lines:    generateLines(fd.Body.List),
	}
}

func generateLines(stmts []ast.Stmt) []line {
	lines := []line{}
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			cx, ok := s.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			switch f := cx.Fun.(type) {
			case *ast.SelectorExpr:
				lines = append(lines, line{s: printLuaFunc(f.X.(*ast.Ident).Name, f.Sel.Name, cx.Args)})
			case *ast.Ident:
				// assume toplevel declared func
				lines = append(lines, line{s: fmt.Sprintf("TODO: %v", f.Name)})
			}
		case *ast.ForStmt:
			if s.Cond == nil {
				fmt.Println("TODO WHILE: ", s)
				continue
			}
			fmt.Println("TODO 3-PART FOR: ", s)
		case *ast.IfStmt:
			cond := "TODO"
			iflines := generateLines(s.Body.List)
			elselines := generateLines(s.Else.(*ast.BlockStmt).List)
			lines = append(lines, line{s: fmt.Sprintf("function () i = mem.condJump(i, %d, %s) end", len(iflines)+2, cond)})
			lines = append(lines, iflines...)
			lines = append(lines, line{s: fmt.Sprintf("function () i = mem.goto(i, %d) end", len(elselines)+1)})
			lines = append(lines, elselines...)
		}
	}
	return lines
}
