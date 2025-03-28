// main пакет для статической проверки кода
package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ErrCheckAnalyzer структура для кастомного обнаружения кода os.Exit в пакете main в функции main
var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for os.Exit in package main and function main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// packageMainName переменная для обнаружения пакета main
	var packageMainName string
	// selName переменная для обнаружения Exit
	var selName string
	// selX переменная для обнаружения os
	var selX string

	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			if v, ok := node.(*ast.File); ok {
				packageMainName = v.Name.String()
			}

			if v, ok := node.(*ast.SelectorExpr); ok {
				selX = fmt.Sprintf("%v", v.X)
				selName = v.Sel.Name
				if packageMainName == "main" && selX == "os" && selName == "Exit" {
					pass.Reportf(v.Pos(), "in package main and function main os.Exit not allowed")
				}
			}

			return true
		})
	}
	return nil, nil
}
