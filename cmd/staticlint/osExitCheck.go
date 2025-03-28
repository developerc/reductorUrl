package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for os.Exit in package main and function main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var packageMainName string //*ast.Ident
	var selName string
	var selX string

	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			if v, ok := node.(*ast.File); ok {
				//fmt.Println(v.Name)
				packageMainName = v.Name.String()
			}

			if v, ok := node.(*ast.SelectorExpr); ok {
				//fmt.Print(n.(*ast.SelectorExpr).X, " ")
				selX = fmt.Sprintf("%v", v.X)
				//fmt.Println(v.Pos(), v.Sel.NamePos, selX)
				//fmt.Print(v.X, " ")
				//fmt.Println(v.Sel.Name)
				selName = v.Sel.Name
				if packageMainName == "main" && selX == "os" && selName == "Exit" {
					//fmt.Print("Yes \n")
					pass.Reportf(v.Pos(), "in package main and function main os.Exit not allowed")
				}
			}

			return true
		})
	}
	return nil, nil
}
