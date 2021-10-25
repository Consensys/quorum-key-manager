package utils

import (
	"strings"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

const usernameTenantSeparator = "|"

func ExtractUsernameAndTenant(sub string) (username, tenant string) {
	if !strings.Contains(sub, usernameTenantSeparator) {
		return sub, ""
	}

	pieces := strings.Split(sub, usernameTenantSeparator)
	return pieces[1], pieces[0]
}

func ExtractPermissionsArr(claims []string) []entities.Permission {
	var permissions []entities.Permission

	for _, claim := range claims {
		if !strings.Contains(claim, ":") {
			// Ignore invalid permissions
			continue
		}

		if strings.Contains(claim, "*") {
			permissions = append(permissions, entities.ListWildcardPermission(claim)...)
		} else {
			permissions = append(permissions, entities.Permission(claim))
		}
	}

	return permissions
}

func ExtractPermissions(claims string) []entities.Permission {
	return ExtractPermissionsArr(strings.Split(claims, " "))
}

func ExtractRoles(roles string) []string {
	return strings.Split(roles, ",")
}
