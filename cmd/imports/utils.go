package imports

// difference returns the elements in `a` that aren't in `b`
// Good algorithm as it's O(n) instead of a naive O(n2)
func difference(a, b []string) []string {
	mb := arrToMap(b)

	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func arrToMap(a []string) map[string]struct{} {
	mb := make(map[string]struct{}, len(a))
	for _, x := range a {
		mb[x] = struct{}{}
	}

	return mb
}
