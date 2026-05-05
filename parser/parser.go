package parser

import (
    "context"

    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/typescript/typescript"
)

func Parse(src []byte) *sitter.Node {
    parser := sitter.NewParser()
    parser.SetLanguage(typescript.GetLanguage())
    tree, err := parser.ParseCtx(context.Background(), nil, src)
    if err != nil {
        panic("parse error: " + err.Error())
    }
    return tree.RootNode()
}
