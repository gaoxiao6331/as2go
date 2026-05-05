package codegen

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

func (g *Generator) genStmt(node *sitter.Node) {
	switch node.Type() {
	case "block":
		g.genBlock(node)
	case "return_statement":
		g.genReturn(node)
	case "if_statement":
		g.genIf(node)
	case "while_statement":
		g.genWhile(node)
	case "do_statement":
		g.genDoWhile(node)
	case "for_statement":
		g.genFor(node)
	case "for_of_statement":
		g.genForOf(node)
	case "break_statement":
		g.writeln("break")
	case "continue_statement":
		g.writeln("continue")
	case "switch_statement":
		g.genSwitch(node)
	case "throw_statement":
		g.genThrow(node)
	case "try_statement":
		panic("as2go: try/catch is not supported yet")
	case "expression_statement":
		g.genExprStmt(node)
	case "lexical_declaration", "variable_declaration":
		g.genLocalVarDecl(node)
	case "empty_statement", "comment":
	default:
		fmt.Printf("// [as2go] unhandled stmt: %s\n", node.Type())
	}
}

func (g *Generator) genBlock(node *sitter.Node) {
	g.writeln("{")
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "{" || child.Type() == "}" {
			continue
		}
		g.genStmt(child)
	}
	g.writeln("}")
}

func (g *Generator) genReturn(node *sitter.Node) {
	hasVal := node.ChildCount() > 1
	if hasVal {
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() != "return" && child.Type() != ";" {
				g.writef("return %s\n", g.genExprStr(child))
				return
			}
		}
	}
	g.writeln("return")
}

func (g *Generator) genIf(node *sitter.Node) {
	cond := g.childByField(node, "condition")
	consequence := g.childByField(node, "consequence")
	alternative := g.childByField(node, "alternative")

	g.writef("if %s ", g.genExprStr(cond))
	g.genBlock(consequence)
	if alternative != nil {
		g.write(" else ")
		alt := alternative
		if alt.Type() == "else_clause" && alt.ChildCount() > 1 {
			inner := alt.Child(1)
			if inner.Type() == "if_statement" {
				g.genIf(inner)
				return
			}
			g.genBlock(inner)
			return
		}
		g.genStmt(alt)
	}
	g.writeln("")
}

func (g *Generator) genWhile(node *sitter.Node) {
	cond := g.childByField(node, "condition")
	body := g.childByField(node, "body")
	g.writef("for %s ", g.genExprStr(cond))
	g.genBlock(body)
	g.writeln("")
}

func (g *Generator) genDoWhile(node *sitter.Node) {
	body := g.childByField(node, "body")
	cond := g.childByField(node, "condition")
	g.writeln("for {")
	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child.Type() == "{" || child.Type() == "}" {
			continue
		}
		g.genStmt(child)
	}
	g.writef("if !(%s) { break }\n", g.genExprStr(cond))
	g.writeln("}")
	g.writeln("")
}

func (g *Generator) genFor(node *sitter.Node) {
	init := g.childByField(node, "initializer")
	cond := g.childByField(node, "condition")
	update := g.childByField(node, "increment")
	body := g.childByField(node, "body")

	initStr := ""
	if init != nil {
		initStr = g.genForInit(init)
	}
	condStr := ""
	if cond != nil {
		condStr = g.genExprStr(cond)
	}
	updateStr := ""
	if update != nil {
		updateStr = g.genExprStr(update)
	}

	g.writef("for %s; %s; %s ", initStr, condStr, updateStr)
	g.genBlock(body)
	g.writeln("")
}

func (g *Generator) genForInit(node *sitter.Node) string {
	switch node.Type() {
	case "lexical_declaration", "variable_declaration":
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "variable_declarator" {
				name := g.text(g.childByField(child, "name"))
				val := g.childByField(child, "value")
				if val != nil {
					return fmt.Sprintf("%s := %s", name, g.genExprStr(val))
				}
				return name
			}
		}
	}
	return g.genExprStr(node)
}

func (g *Generator) genForOf(node *sitter.Node) {
	left := g.childByField(node, "left")
	right := g.childByField(node, "right")
	body := g.childByField(node, "body")

	varName := g.text(left)
	if left.Type() == "lexical_declaration" || left.Type() == "variable_declaration" {
		for i := 0; i < int(left.ChildCount()); i++ {
			child := left.Child(i)
			if child.Type() == "variable_declarator" {
				varName = g.text(g.childByField(child, "name"))
			}
		}
	}

	g.writef("for _, %s := range %s ", varName, g.genExprStr(right))
	g.genBlock(body)
	g.writeln("")
}

func (g *Generator) genSwitch(node *sitter.Node) {
	val := g.childByField(node, "value")
	body := g.childByField(node, "body")

	g.writef("switch %s {\n", g.genExprStr(val))
	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child.Type() == "switch_case" {
			caseVal := g.childByField(child, "value")
			g.writef("case %s:\n", g.genExprStr(caseVal))
			for j := 0; j < int(child.ChildCount()); j++ {
				s := child.Child(j)
				if s.Type() != "case" && s.Type() != ":" && s != caseVal {
					g.genStmt(s)
				}
			}
		} else if child.Type() == "switch_default" {
			g.writeln("default:")
			for j := 0; j < int(child.ChildCount()); j++ {
				s := child.Child(j)
				if s.Type() != "default" && s.Type() != ":" {
					g.genStmt(s)
				}
			}
		}
	}
	g.writeln("}")
	g.writeln("")
}

func (g *Generator) genThrow(node *sitter.Node) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() != "throw" && child.Type() != ";" {
			g.writef("panic(%s)\n", g.genExprStr(child))
			return
		}
	}
}

func (g *Generator) genExprStmt(node *sitter.Node) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() != ";" {
			g.writeln(g.genExprStr(child))
		}
	}
}

func (g *Generator) genLocalVarDecl(node *sitter.Node) {
	isConst := node.Child(0) != nil && g.text(node.Child(0)) == "const"
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "variable_declarator" {
			name := g.text(g.childByField(child, "name"))
			typ := g.childByField(child, "type")
			val := g.childByField(child, "value")

			if isConst {
				if val != nil {
					g.writef("const %s = %s\n", name, g.genExprStr(val))
				}
			} else if val != nil {
				if typ != nil {
					g.writef("%s := %s(%s)\n", name, g.genType(typ), g.genExprStr(val))
				} else {
					g.writef("%s := %s\n", name, g.genExprStr(val))
				}
			} else {
				goType := ""
				if typ != nil {
					goType = g.genType(typ)
				}
				g.writef("var %s %s\n", name, goType)
			}
		}
	}
}
