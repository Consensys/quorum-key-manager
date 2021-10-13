package imports

// difference returns the elements in `a` that aren't in `b`
// Good algorithm as it's O(n) instead of a naive O(n2)
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func contains(val string, arr []string) bool {
	for _, x := range arr {
		if val == x {
			return true
		}
	}

	return false
}
