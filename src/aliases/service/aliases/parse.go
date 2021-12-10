package aliases

import "strings"

func (s *Aliases) Parse(alias string) (regName, aliasKey string, isAlias bool) {
	chunks := strings.Split(alias, ":")

	if len(chunks) != 2 {
		return "", "", false
	}

	return chunks[0], chunks[1], true
}
