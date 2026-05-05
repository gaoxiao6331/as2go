package main

import (
    "fmt"
    "go/format"
    "os"

    "as2go/codegen"
    "as2go/parser"
    "as2go/preprocess"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: as2go <input.ts>")
        os.Exit(1)
    }

    src, err := os.ReadFile(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    processed := preprocess.Process(src)
    tree := parser.Parse(processed)

    gen := codegen.New()
    goSrc := gen.Generate(tree, processed)

    formatted, err := format.Source([]byte(goSrc))
    if err != nil {
        fmt.Println(goSrc)
        fmt.Fprintln(os.Stderr, "format error:", err)
        os.Exit(1)
    }

    fmt.Print(string(formatted))
}
