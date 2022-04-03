package common

// FirstSet returns first non-nil argument passed
func FirstSet(s ...*string) *string {
	for _, elem := range s {
		if elem != nil {
			return elem
		}
	}
	return nil
}
