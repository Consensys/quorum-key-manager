package jose

import (
	"time"
)

type Config struct {
	IssuerURL           string
	CacheTTL            time.Duration
	Audience            []string
	CustomClaimPath     string
	PermissionClaimPath string
	RolesClaimPath      string
}

func NewConfig(issuerURL string, audience []string, customClaimPath, permissionClaimPath, rolesClaimPath string, cacheTTL time.Duration) *Config {
	return &Config{
		IssuerURL:           issuerURL,
		CacheTTL:            cacheTTL,
		Audience:            audience,
		CustomClaimPath:     customClaimPath,
		PermissionClaimPath: permissionClaimPath,
		RolesClaimPath:      rolesClaimPath,
	}
}
