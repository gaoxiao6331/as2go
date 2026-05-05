package codegen

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

type Generator struct {
	src     []byte
	imports map[string]bool
	runtime map[string]bool
	buf     strings.Builder
	indent  int
}

func New() *Generator {
	return &Generator{
		imports: make(map[string]bool),
		runtime: make(map[string]bool),
	}
}

func (g *Generator) Generate(root *sitter.Node, src []byte) string {
	g.src = src

	g.writeln("package main")
	g.writeln("")

	var body strings.Builder
	bodyGen := &Generator{src: src, imports: g.imports, runtime: g.runtime}
	bodyGen.visitNode(root)
	body.WriteString(bodyGen.buf.String())

	if len(g.imports) > 0 {
		g.writeln("import (")
		for pkg := range g.imports {
			g.writef("    \"%s\"", pkg)
			g.writeln("")
		}
		g.writeln(")")
		g.writeln("")
	}

	if len(g.runtime) > 0 {
		g.writeln(strings.TrimSpace(runtimeFuncs))
		g.writeln("")
	}

	g.buf.WriteString(body.String())
	return g.buf.String()
}

func (g *Generator) visitNode(node *sitter.Node) {
	switch node.Type() {
	case "program":
		g.visitChildren(node)

	case "function_declaration":
		g.genFuncDecl(node)
	case "class_declaration":
		g.genClassDecl(node)
	case "enum_declaration":
		g.genEnumDecl(node)
	case "interface_declaration":
		g.genInterfaceDecl(node)
	case "lexical_declaration", "variable_declaration":
		g.genVarDecl(node)
	case "type_alias_declaration":
		g.genTypeAlias(node)
	case "import_statement":
		g.genImport(node)
	case "export_statement":
		g.genExport(node)
	case "expression_statement":
		g.genExprStmt(node)

	case "comment", "empty_statement":

	default:
		fmt.Printf("// [as2go] unhandled node type at top level: %s\n", node.Type())
	}
}

func (g *Generator) visitChildren(node *sitter.Node) {
	for i := 0; i < int(node.ChildCount()); i++ {
		g.visitNode(node.Child(i))
	}
}

func (g *Generator) text(node *sitter.Node) string {
	return node.Content(g.src)
}

func (g *Generator) write(s string) {
	g.buf.WriteString(s)
}

func (g *Generator) writeln(s string) {
	g.buf.WriteString(s)
	g.buf.WriteByte('\n')
}

func (g *Generator) writef(format string, args ...any) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) addImport(pkg string) {
	if g.imports == nil {
		g.imports = make(map[string]bool)
	}
	g.imports[pkg] = true
}

func (g *Generator) requireRuntime(name string) {
	if g.runtime == nil {
		g.runtime = make(map[string]bool)
	}
	g.runtime[name] = true
}
