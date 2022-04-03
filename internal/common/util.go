package common

import "flag"

// IsFlagPassed checks if the specified flag was passed via the command line
func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// CopyIfNotNil copies src to dst if src is not nil
func CopyIfNotNil(dst, src *string) {
	if src != nil {
		*dst = *src
	}
}
