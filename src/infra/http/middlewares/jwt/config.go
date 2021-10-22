package jwt

import (
	"time"
)

type Config struct {
	IssuerURL string
	CacheTTL  time.Duration
	Claims    ClaimsConfig
}

type ClaimsConfig struct {
	Subject string
	Scope   string
	Roles   string
}

func NewConfig(issuerURL, subject, scope, roles string, cacheTTL time.Duration) *Config {
	return &Config{
		IssuerURL: issuerURL,
		CacheTTL:  cacheTTL,
		Claims: ClaimsConfig{
			Subject: subject,
			Scope:   scope,
			Roles:   roles,
		},
	}
}
