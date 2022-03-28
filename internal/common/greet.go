package common

import "fmt"

func na(a string) string {
	if a == "" {
		return "N/A"
	}
	return a
}

// PrintBuildInfo prins build version, date and commit
func PrintBuildInfo(bv, bd, bc string) {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		na(bv), na(bd), na(bc))
}
