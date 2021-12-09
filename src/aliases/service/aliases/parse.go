package aliases

func (i *Aliases) Parse(alias string) (regName, aliasKey string, isAlias bool) {
	submatches := i.regex.FindStringSubmatch(alias)
	if len(submatches) < 3 {
		return "", "", false
	}

	regName = submatches[1]
	aliasKey = submatches[2]

	return regName, aliasKey, true
}
