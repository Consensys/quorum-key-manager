package aliases

import "strings"

func (s *Aliases) Parse(alias string) (regName, aliasKey string, isAlias bool) {
	if strings.HasPrefix(alias, "{{") && strings.HasSuffix(alias, "}}") {
		chunks := strings.Split(alias[2:len(alias)-2], ":")

		if len(chunks) != 2 {
			return "", "", false
		}

		return chunks[0], chunks[1], true
	}

	return "", "", false
}
