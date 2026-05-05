package codegen

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

func (g *Generator) genFuncDecl(node *sitter.Node) {
	name := g.childByField(node, "name")
	params := g.childByField(node, "parameters")
	retType := g.childByField(node, "return_type")
	body := g.childByField(node, "body")

	g.write("func ")
	g.write(g.text(name))
	g.write("(")
	g.genParams(params)
	g.write(")")

	if retType != nil {
		rt := g.genType(retType)
		if rt != "" {
			g.write(" " + rt)
		}
	}

	g.write(" ")
	g.genBlock(body)
	g.writeln("")
}

func (g *Generator) genMethodDecl(node *sitter.Node, receiver string) {
	name := g.childByField(node, "name")
	params := g.childByField(node, "parameters")
	retType := g.childByField(node, "return_type")
	body := g.childByField(node, "body")

	g.writef("func (self *%s) %s(", receiver, g.text(name))
	g.genParams(params)
	g.write(")")

	if retType != nil {
		rt := g.genType(retType)
		if rt != "" {
			g.write(" " + rt)
		}
	}

	g.write(" ")
	g.genBlock(body)
	g.writeln("")
}

func (g *Generator) genClassDecl(node *sitter.Node) {
	name := g.text(g.childByField(node, "name"))
	heritage := g.childByField(node, "class_heritage")
	body := g.childByField(node, "body")

	if heritage != nil {
		panic(fmt.Sprintf("as2go: class inheritance (extends) is not supported yet, class: %s", name))
	}

	g.writef("type %s struct {\n", name)
	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child.Type() == "public_field_definition" || child.Type() == "field_definition" {
			fieldName := g.childByField(child, "name")
			fieldType := g.childByField(child, "type")
			goType := ""
			if fieldType != nil {
				goType = g.genType(fieldType)
			}
			fn := g.text(fieldName)
			g.writef("    %s %s\n", capitalize(fn), goType)
		}
	}
	g.writeln("}")
	g.writeln("")

	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child.Type() == "method_definition" {
			g.genMethodDecl(child, name)
		}
	}
}

func (g *Generator) genEnumDecl(node *sitter.Node) {
	name := g.text(g.childByField(node, "name"))
	body := g.childByField(node, "body")

	g.writef("type %s int32\n", name)
	g.writeln("const (")
	first := true
	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child.Type() == "enum_assignment" || child.Type() == "property_identifier" {
			var memberName string
			var memberVal *sitter.Node
			if child.Type() == "enum_assignment" {
				memberName = g.text(g.childByField(child, "name"))
				memberVal = g.childByField(child, "value")
			} else {
				memberName = g.text(child)
			}
			if first {
				if memberVal != nil {
					g.writef("    %s_%s %s = %s\n", name, memberName, name, g.genExprStr(memberVal))
				} else {
					g.writef("    %s_%s %s = iota\n", name, memberName, name)
				}
				first = false
			} else {
				if memberVal != nil {
					g.writef("    %s_%s = %s\n", name, memberName, g.genExprStr(memberVal))
				} else {
					g.writef("    %s_%s\n", name, memberName)
				}
			}
		}
	}
	g.writeln(")")
	g.writeln("")
}

func (g *Generator) genInterfaceDecl(node *sitter.Node) {
	name := g.text(g.childByField(node, "name"))
	body := g.childByField(node, "body")

	g.writef("type %s interface {\n", name)
	for i := 0; i < int(body.ChildCount()); i++ {
		child := body.Child(i)
		if child.Type() == "method_signature" {
			methodName := g.childByField(child, "name")
			params := g.childByField(child, "parameters")
			retType := g.childByField(child, "return_type")
			g.writef("    %s(", g.text(methodName))
			g.genParams(params)
			g.write(")")
			if retType != nil {
				rt := g.genType(retType)
				if rt != "" {
					g.write(" " + rt)
				}
			}
			g.writeln("")
		}
	}
	g.writeln("}")
	g.writeln("")
}

func (g *Generator) genVarDecl(node *sitter.Node) {
	kind := node.Child(0)
	isConst := kind != nil && g.text(kind) == "const"

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "variable_declarator" {
			varName := g.text(g.childByField(child, "name"))
			varType := g.childByField(child, "type")
			varVal := g.childByField(child, "value")
			if isConst {
				g.write("const ")
			} else {
				g.write("var ")
			}
			g.write(varName)
			if varType != nil {
				g.write(" " + g.genType(varType))
			}
			if varVal != nil {
				g.write(" = " + g.genExprStr(varVal))
			}
			g.writeln("")
		}
	}
}

func (g *Generator) genTypeAlias(node *sitter.Node) {
	name := g.text(g.childByField(node, "name"))
	typeVal := g.childByField(node, "value")
	g.writef("type %s = %s\n\n", name, g.genType(typeVal))
}

func (g *Generator) genImport(node *sitter.Node) {
	fmt.Printf("// [as2go] import statement skipped: %s\n", g.text(node))
}

func (g *Generator) genExport(node *sitter.Node) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		switch child.Type() {
		case "function_declaration":
			g.genFuncDecl(child)
		case "class_declaration":
			g.genClassDecl(child)
		case "enum_declaration":
			g.genEnumDecl(child)
		case "interface_declaration":
			g.genInterfaceDecl(child)
		case "lexical_declaration", "variable_declaration":
			g.genVarDecl(child)
		}
	}
}

func (g *Generator) genParams(node *sitter.Node) {
	if node == nil {
		return
	}
	first := true
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "required_parameter" || child.Type() == "optional_parameter" {
			if !first {
				g.write(", ")
			}
			first = false
			paramName := g.childByField(child, "pattern")
			paramType := g.childByField(child, "type")
			g.write(g.text(paramName))
			if paramType != nil {
				g.write(" " + g.genType(paramType))
			}
		}
	}
}

func (g *Generator) genType(node *sitter.Node) string {
	if node == nil {
		return ""
	}
	switch node.Type() {
	case "type_annotation":
		if node.ChildCount() > 1 {
			return g.genType(node.Child(1))
		}
		return ""
	case "predefined_type":
		return mapType(g.text(node))
	case "type_identifier":
		return mapType(g.text(node))
	case "generic_type":
		baseName := g.text(g.childByField(node, "name"))
		args := g.childByField(node, "type_arguments")
		return g.genGenericType(baseName, args)
	case "array_type":
		elem := g.genType(node.Child(0))
		return "[]" + elem
	case "union_type":
		panic("as2go: union types are not supported (use a concrete type instead)")
	case "void_type":
		return ""
	default:
		return g.text(node)
	}
}

func (g *Generator) genGenericType(name string, args *sitter.Node) string {
	typeArgs := g.collectTypeArgs(args)
	switch name {
	case "Array":
		if len(typeArgs) != 1 {
			panic("as2go: Array<T> requires exactly one type argument")
		}
		return "[]" + typeArgs[0]
	case "StaticArray":
		if len(typeArgs) != 2 {
			panic("as2go: StaticArray<T,N> requires two type arguments")
		}
		return fmt.Sprintf("[%s]%s", typeArgs[1], typeArgs[0])
	case "Map":
		if len(typeArgs) != 2 {
			panic("as2go: Map<K,V> requires two type arguments")
		}
		return fmt.Sprintf("map[%s]%s", typeArgs[0], typeArgs[1])
	case "Set":
		if len(typeArgs) != 1 {
			panic("as2go: Set<T> requires one type argument")
		}
		return fmt.Sprintf("map[%s]struct{}", typeArgs[0])
	default:
		panic(fmt.Sprintf("as2go: generic type %s is not supported yet", name))
	}
}

func (g *Generator) collectTypeArgs(node *sitter.Node) []string {
	if node == nil {
		return nil
	}
	var args []string
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		t := child.Type()
		if t == "," || t == "<" || t == ">" {
			continue
		}
		args = append(args, g.genType(child))
	}
	return args
}

func (g *Generator) childByField(node *sitter.Node, field string) *sitter.Node {
	return node.ChildByFieldName(field)
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
