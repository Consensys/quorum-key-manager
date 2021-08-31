package utils

import (
	"strings"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const usernameTenantSeparator = "|"

func ExtractUsernameAndTenant(sub string) (username, tenant string) {
	if !strings.Contains(sub, usernameTenantSeparator) {
		return sub, ""
	}

	pieces := strings.Split(sub, usernameTenantSeparator)
	return pieces[1], pieces[0]
}

func ExtractPermissions(claims []string) []types.Permission {
	var permissions []types.Permission

	for _, claim := range claims {
		switch {
		case strings.Contains(claim, " "):
			subPermissions := ExtractPermissions(strings.Split(claim, " "))
			permissions = append(permissions, subPermissions...)

		case strings.Contains(claim, ":"):
			if strings.Contains(claim, "*") {
				permissions = append(permissions, types.ListWildcardPermission(claim)...)
			} else {
				permissions = append(permissions, types.Permission(claim))
			}
		}
	}

	return permissions
}

func ExtractRoles(roles string) []string {
	return strings.Split(roles, ",")
}
