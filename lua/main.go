package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"math/rand"
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

// TODO can include variables for goto etc (?)
type line struct {
	s string
}

var gensymFunc func() string = gensym

func gensym() string {
	return fmt.Sprintf("gensym%d", rand.Int())
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
			fmt.Println(generateProgram(f))
		}
	}
}

// a program is an exported func of one parameter, which is a turtle
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

// TODO: this should generate version of program with mem lib
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

// generates a turtle program in lua that uses memlib so it keeps state over server downtime
// and in general the chunk being offloaded in minecraft
func generateProgram(fd *ast.FuncDecl) string {
	g := &generator{}
	name := strings.ToLower(fd.Name.Name)
	// header
	g.Writef(`-- comment
local state = mem.startFromMemory(%q)
if not state then
    state = {i=0}
end
local i = state.i

local stop = false -- used to communicate key Q pressed

action = {
`, name)
	// list of statements that should each be atomic wrt computercraft and minecraft ticks
	block := generateFirstPass(fd)
	for _, line := range block.lines {
		g.Writef("    " + line.s + ",\n")
	}
	// footer incl main call
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
			lines = append(lines, generateExpr(s.X))
		case *ast.ForStmt:
			body := generateLines(s.Body.List)
			if s.Cond == nil {
				lines = append(lines, body...)
				lines = append(lines, line{s: fmt.Sprintf("function () i = mem.goto(i, %d) end", -(len(body) + 1))})
				continue
			}
			cond := negate(s.Cond)
			c := generateExpr(cond)
			if s.Init == nil {
				lines = append(lines, line{s: fmt.Sprintf("function () i = mem.condJump(i, %d, %s) end", len(body)+2, c.s)})
				lines = append(lines, body...)
				lines = append(lines, line{s: fmt.Sprintf("function () i = mem.goto(i, %d) end", -(len(body) + 1))})
				continue
			}
			// here we introduce a new variable into the global state!
			// TODO: assumes init is var:=0 and incr is ++
			loopvar := gensymFunc()
			lines = append(lines, line{s: fmt.Sprintf("function () if state.%s == nil then state.%s = 0 end i = mem.condJump(i, %d, state.%s) end", loopvar, loopvar, len(body)+2, c.s)})
			body[len(body)-1].s = "function () " + body[len(body)-1].s + fmt.Sprintf("; state.%s = state.%s+1", loopvar, loopvar)
			lines = append(lines, body...)
			lines = append(lines, line{s: fmt.Sprintf("function () i = mem.goto(i, %d) end", -(len(body) + 1))})
		case *ast.IfStmt:
			cond := negate(s.Cond)
			c := generateExpr(cond)
			iflines := generateLines(s.Body.List)
			elselines := generateLines(s.Else.(*ast.BlockStmt).List)
			lines = append(lines, line{s: fmt.Sprintf("function () i = mem.condJump(i, %d, %s) end", len(iflines)+2, c.s)})
			lines = append(lines, iflines...)
			lines = append(lines, line{s: fmt.Sprintf("function () i = mem.goto(i, %d) end", len(elselines)+1)})
			lines = append(lines, elselines...)
		default:
			fmt.Printf("unexpected ast type %T for stmt\n", s)
		}
	}
	return lines
}

func generateExpr(s ast.Expr) line {
	switch t := s.(type) {
	case *ast.CallExpr:
		return generateCallExpr(t)
	case *ast.BinaryExpr:
		return generateBinaryExpr(t)
	case *ast.UnaryExpr:
		return generateUnaryExpr(t)
	case *ast.Ident:
		return line{s: t.Name}
	case *ast.BasicLit:
		return line{s: t.Value}
	default:
		fmt.Printf("unexpected ast type %T for expr\n", t)
		return line{}
	}
}

func generateCallExpr(s *ast.CallExpr) line {
	switch f := s.Fun.(type) {
	case *ast.SelectorExpr:
		return line{s: printLuaFunc(f.X.(*ast.Ident).Name, f.Sel.Name, s.Args)}
	case *ast.Ident:
		// assume toplevel declared func
		return line{s: fmt.Sprintf("TODO: %v", f.Name)}
	default:
		fmt.Printf("unexpected ast type %T for callexpr func\n", f)
		return line{}
	}
}

func generateBinaryExpr(s *ast.BinaryExpr) line {
	switch s.Op {
	case token.LSS:
		return line{s: generateExpr(s.X).s + " < " + generateExpr(s.Y).s}
	case token.LEQ:
		return line{s: generateExpr(s.X).s + " <= " + generateExpr(s.Y).s}
	case token.GTR:
		return line{s: generateExpr(s.X).s + " < " + generateExpr(s.Y).s}
	case token.GEQ:
		return line{s: generateExpr(s.X).s + " >= " + generateExpr(s.Y).s}
	default:
		fmt.Printf("unsupported operator %#v for binaryexpr\n", s.Op.String())
		return line{}
	}
}

func generateUnaryExpr(s *ast.UnaryExpr) line {
	switch s.Op {
	case token.NOT:
		return line{s: "not " + generateExpr(s.X).s}
	default:
		fmt.Printf("unsupported operator %#v for unaryexpr\n", s.Op.String())
		return line{}
	}
}

func negate(s ast.Expr) ast.Expr {
	switch t := s.(type) {
	case *ast.UnaryExpr:
		switch t.Op {
		case token.NOT:
			return t.X
		default:
			fmt.Printf("unsupported operator %#v for negate unaryexpr\n", t.Op.String())
			return t
		}
	case *ast.BinaryExpr:
		switch t.Op {
		case token.LSS:
			return &ast.BinaryExpr{X: t.X, Y: t.Y, Op: token.GEQ}
		case token.LEQ:
			return &ast.BinaryExpr{X: t.X, Y: t.Y, Op: token.GTR}
		case token.GTR:
			return &ast.BinaryExpr{X: t.X, Y: t.Y, Op: token.LEQ}
		case token.GEQ:
			return &ast.BinaryExpr{X: t.X, Y: t.Y, Op: token.LSS}
		default:
			fmt.Printf("unsupported operator %#v for negate binaryexpr\n", t.Op.String())
			return t
		}
	default:
		return &ast.UnaryExpr{Op: token.NOT, X: t}
	}
}
