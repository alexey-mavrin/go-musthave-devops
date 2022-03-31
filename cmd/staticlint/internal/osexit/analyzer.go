package osexit

import (
	"go/ast"
	"go/printer"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// OSExitAnalyzer is os.Exit-in-main analyzer
var OSExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.Exit in main function",
	Run:  runExitCheck,
}

func runExitCheck(pass *analysis.Pass) (interface{}, error) {
	var b strings.Builder
	var currentPackage string
	var currentFunction string
	for _, file := range pass.Files {
		currentPackage = file.Name.Name
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr:
				pos := x.Fun.Pos()
				printer.Fprint(&b, pass.Fset, x.Fun)
				if b.String() == "os.Exit" &&
					currentPackage == "main" &&
					currentFunction == "main" {
					pass.Reportf(pos, "os.Exit call in main.main")
				}
			case *ast.FuncDecl:
				currentFunction = x.Name.Name
			}
			b.Reset()
			return true
		})
	}
	return nil, nil
}
