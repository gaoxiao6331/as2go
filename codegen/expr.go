package codegen

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

func (g *Generator) genExprStr(node *sitter.Node) string {
	if node == nil {
		return ""
	}
	switch node.Type() {
	case "identifier":
		name := g.text(node)
		if name == "Math" {
			g.addImport("math")
			return "math"
		}
		if name == "abort" {
			g.requireRuntime("abort")
			return "__as_abort"
		}
		return name
	case "this":
		return "self"
	case "super":
		panic("as2go: super is not supported yet (class inheritance not implemented)")
	case "true":
		return "true"
	case "false":
		return "false"
	case "null":
		return "nil"
	case "number":
		return g.text(node)
	case "string":
		return g.text(node)
	case "parenthesized_expression":
		inner := node.Child(1)
		return "(" + g.genExprStr(inner) + ")"
	case "binary_expression":
		return g.genBinary(node)
	case "unary_expression":
		op := g.text(node.Child(0))
		operand := node.Child(1)
		return op + g.genExprStr(operand)
	case "update_expression":
		if node.Child(0).Type() == "++" || node.Child(0).Type() == "--" {
			return g.text(node.Child(0)) + g.genExprStr(node.Child(1))
		}
		return g.genExprStr(node.Child(0)) + g.text(node.Child(1))
	case "assignment_expression":
		left := g.childByField(node, "left")
		right := g.childByField(node, "right")
		op := ""
		for i := 0; i < int(node.ChildCount()); i++ {
			t := node.Child(i).Type()
			if strings.HasSuffix(t, "=") && t != "==" && t != "!=" && t != "<=" && t != ">=" {
				op = t
				break
			}
		}
		return fmt.Sprintf("%s %s %s", g.genExprStr(left), op, g.genExprStr(right))
	case "call_expression":
		return g.genCall(node)
	case "new_expression":
		return g.genNew(node)
	case "member_expression":
		obj := g.childByField(node, "object")
		prop := g.childByField(node, "property")
		base := g.genExprStr(obj)
		propName := g.text(prop)

		if propName == "length" {
			return fmt.Sprintf("len(%s)", base)
		}

		if obj.Type() == "this" {
			return base + "." + capitalize(propName)
		}

		if base == "math" {
			return base + "." + capitalize(propName)
		}

		return base + "." + propName
	case "subscript_expression":
		obj := g.childByField(node, "object")
		idx := g.childByField(node, "index")
		return fmt.Sprintf("%s[%s]", g.genExprStr(obj), g.genExprStr(idx))
	case "augmented_assignment_expression":
		left := g.childByField(node, "left")
		right := g.childByField(node, "right")
		opNode := node.Child(1)
		return fmt.Sprintf("%s %s %s", g.genExprStr(left), g.text(opNode), g.genExprStr(right))
	case "ternary_expression":
		panic("as2go: ternary expression is not supported yet (extract to if statement)")
	case "type_assertion_expression", "as_expression":
		val := node.Child(0)
		typ := node.Child(2)
		return fmt.Sprintf("%s(%s)", g.genType(typ), g.genExprStr(val))
	case "instanceof_expression":
		val := node.Child(0)
		typ := node.Child(2)
		return fmt.Sprintf("func() bool { _, ok := %s.(%s); return ok }()", g.genExprStr(val), g.genType(typ))
	case "array":
		return g.genArrayLiteral(node)
	case "object":
		panic("as2go: object literal is not supported yet (use a struct instead)")
	case "arrow_function", "function_expression":
		return g.genFuncLiteral(node)
	case "await_expression":
		panic("as2go: async/await is not supported yet")
	case "non_null_expression":
		return g.genExprStr(node.Child(0))
	case "expression_statement":
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() != ";" {
				return g.genExprStr(child)
			}
		}
		return ""
	case "comma_expression":
		panic("as2go: comma expression is not supported yet (split into separate statements)")
	default:
		fmt.Printf("// [as2go] unhandled expr: %s\n", node.Type())
		return g.text(node)
	}
}

func (g *Generator) genBinary(node *sitter.Node) string {
	left := g.childByField(node, "left")
	right := g.childByField(node, "right")
	op := ""
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		t := child.Type()
		if child != left && child != right {
			op = t
			break
		}
	}
	switch op {
	case "===":
		op = "=="
	case "!==":
		op = "!="
	}
	return fmt.Sprintf("%s %s %s", g.genExprStr(left), op, g.genExprStr(right))
}

func (g *Generator) genCall(node *sitter.Node) string {
	fn := g.childByField(node, "function")
	args := g.childByField(node, "arguments")

	fnStr := g.genExprStr(fn)

	switch fnStr {
	case "load":
		panic("as2go: load<T>() memory intrinsic is not supported yet")
	case "store":
		panic("as2go: store<T>() memory intrinsic is not supported yet")
	case "changetype":
		panic("as2go: changetype<T>() is not supported yet")
	case "memory.grow", "memory.size", "memory.fill", "memory.copy":
		panic(fmt.Sprintf("as2go: %s() memory intrinsic is not supported yet", fnStr))
	}

	var argStrs []string
	if args != nil {
		for i := 0; i < int(args.ChildCount()); i++ {
			child := args.Child(i)
			if child.Type() != "," && child.Type() != "(" && child.Type() != ")" {
				argStrs = append(argStrs, g.genExprStr(child))
			}
		}
	}

	return fmt.Sprintf("%s(%s)", fnStr, strings.Join(argStrs, ", "))
}

func (g *Generator) genNew(node *sitter.Node) string {
	constructor := g.childByField(node, "constructor")
	args := g.childByField(node, "arguments")

	typeName := g.text(constructor)
	var argStrs []string
	if args != nil {
		for i := 0; i < int(args.ChildCount()); i++ {
			child := args.Child(i)
			if child.Type() != "," && child.Type() != "(" && child.Type() != ")" {
				argStrs = append(argStrs, g.genExprStr(child))
			}
		}
	}

	if len(argStrs) == 0 {
		return fmt.Sprintf("&%s{}", typeName)
	}
	return fmt.Sprintf("New%s(%s)", typeName, strings.Join(argStrs, ", "))
}

func (g *Generator) genArrayLiteral(node *sitter.Node) string {
	var elems []string
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() != "," && child.Type() != "[" && child.Type() != "]" {
			elems = append(elems, g.genExprStr(child))
		}
	}
	return fmt.Sprintf("[]any{%s}", strings.Join(elems, ", "))
}

func (g *Generator) genFuncLiteral(node *sitter.Node) string {
	params := g.childByField(node, "parameters")
	retType := g.childByField(node, "return_type")
	body := g.childByField(node, "body")

	var sb strings.Builder
	sb.WriteString("func(")
	inner := &Generator{src: g.src, imports: g.imports, runtime: g.runtime}
	inner.genParams(params)
	sb.WriteString(inner.buf.String())
	sb.WriteString(")")

	if retType != nil {
		rt := g.genType(retType)
		if rt != "" {
			sb.WriteString(" " + rt)
		}
	}
	sb.WriteString(" ")

	bodyGen := &Generator{src: g.src, imports: g.imports, runtime: g.runtime}
	if body.Type() == "statement_block" || body.Type() == "block" {
		bodyGen.genBlock(body)
	} else {
		bodyGen.writef("{ return %s }", g.genExprStr(body))
	}
	sb.WriteString(bodyGen.buf.String())

	return sb.String()
}
