package common

// ShortString makes hashes short for a limited column size
func ShortString(s string, tailLength int) string {
	runes := []rune(s)
	if len(runes)/2 > tailLength {
		return string(runes[:tailLength]) + "..." + string(runes[len(runes)-tailLength:])
	}
	return s
}
