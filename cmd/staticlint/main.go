//go:build ignore
// +build ignore

package main

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/alexey-mavrin/go-musthave-devops/cmd/staticlint/internal/osexit"
)

func addPasses(mychecks []*analysis.Analyzer) []*analysis.Analyzer {
	mychecks = append(mychecks, asmdecl.Analyzer)
	mychecks = append(mychecks, assign.Analyzer)
	mychecks = append(mychecks, atomic.Analyzer)
	mychecks = append(mychecks, atomicalign.Analyzer)
	mychecks = append(mychecks, bools.Analyzer)
	mychecks = append(mychecks, buildssa.Analyzer)
	mychecks = append(mychecks, buildtag.Analyzer)
	mychecks = append(mychecks, cgocall.Analyzer)
	mychecks = append(mychecks, composite.Analyzer)
	mychecks = append(mychecks, copylock.Analyzer)
	mychecks = append(mychecks, ctrlflow.Analyzer)
	mychecks = append(mychecks, deepequalerrors.Analyzer)
	mychecks = append(mychecks, errorsas.Analyzer)
	mychecks = append(mychecks, fieldalignment.Analyzer)
	mychecks = append(mychecks, findcall.Analyzer)
	mychecks = append(mychecks, framepointer.Analyzer)
	mychecks = append(mychecks, httpresponse.Analyzer)
	mychecks = append(mychecks, ifaceassert.Analyzer)
	mychecks = append(mychecks, inspect.Analyzer)
	mychecks = append(mychecks, loopclosure.Analyzer)
	mychecks = append(mychecks, lostcancel.Analyzer)
	mychecks = append(mychecks, nilfunc.Analyzer)
	mychecks = append(mychecks, nilness.Analyzer)
	mychecks = append(mychecks, pkgfact.Analyzer)
	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, reflectvaluecompare.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, shift.Analyzer)
	mychecks = append(mychecks, sigchanyzer.Analyzer)
	mychecks = append(mychecks, sortslice.Analyzer)
	mychecks = append(mychecks, stdmethods.Analyzer)
	mychecks = append(mychecks, stringintconv.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)
	mychecks = append(mychecks, testinggoroutine.Analyzer)
	mychecks = append(mychecks, tests.Analyzer)
	mychecks = append(mychecks, unmarshal.Analyzer)
	mychecks = append(mychecks, unreachable.Analyzer)
	mychecks = append(mychecks, unsafeptr.Analyzer)
	mychecks = append(mychecks, unusedresult.Analyzer)
	mychecks = append(mychecks, unusedwrite.Analyzer)
	mychecks = append(mychecks, usesgenerics.Analyzer)

	return mychecks
}

func addSAChecks(mychecks []*analysis.Analyzer) []*analysis.Analyzer {
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	for _, v := range simple.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "S1000") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	for _, v := range stylecheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "ST1000") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	for _, v := range quickfix.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "QF1001") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	return mychecks
}

func main() {
	var mychecks []*analysis.Analyzer
	mychecks = addPasses(mychecks)
	mychecks = addSAChecks(mychecks)
	mychecks = append(mychecks, osexit.Analyzer)

	multichecker.Main(
		mychecks...,
	)
}
