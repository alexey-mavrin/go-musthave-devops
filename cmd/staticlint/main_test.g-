package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
)

func Test_AddChecker(t *testing.T) {
	type args struct {
		mychecks []*analysis.Analyzer
	}
	tests := []struct {
		name     string
		function func([]*analysis.Analyzer) []*analysis.Analyzer
	}{
		{
			name:     "addSAChecks add checks",
			function: addSAChecks,
		},
		{
			name:     "addPasses add checks",
			function: addPasses,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var arg []*analysis.Analyzer
			got := tt.function(arg)
			// check that function actualy adds checkers
			assert.Positive(t, len(got))
		})
	}
}
