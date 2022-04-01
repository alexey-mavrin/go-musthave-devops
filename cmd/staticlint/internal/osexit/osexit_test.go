package osexit_test

import (
	_ "embed"
	"fmt"
	"log"
	"path"
	"testing"

	"github.com/alexey-mavrin/go-musthave-devops/cmd/staticlint/internal/osexit"
	"golang.org/x/tools/go/analysis/analysistest"
)

//go:embed testdata/main
var mainFile string

//go:embed testdata/t
var tFile string

func Test(t *testing.T) {
	var files = make(map[string]string)

	files["main.go"] = mainFile
	files["t/t.go"] = tFile

	dir, _, err := analysistest.WriteFiles(files)
	if err != nil {
		log.Fatal(err)
	}
	res := analysistest.Run(t, path.Join(dir, "src"), osexit.Analyzer)

	for r := range res {
		fmt.Print(r)
	}

	// clean()
}
