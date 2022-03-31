package common

import "fmt"

// PrintBuildInfo prins build version, date and commit
func PrintBuildInfo(bv, bd, bc string) {
	if bv == "" {
		bv = "N/A"
	}

	if bd == "" {
		bd = "N/A"
	}

	if bc == "" {
		bc = "N/A"
	}

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		bv, bd, bc)
}
