package jose

import (
	"time"
)

type Config struct {
	IssuerURL       string
	CacheTTL        time.Duration
	Audience        []string
	CustomClaimPath string
}

func NewConfig(issuerURL string, audience []string, customClaimPath string, cacheTTL time.Duration) *Config {
	return &Config{
		IssuerURL:       issuerURL,
		CacheTTL:        cacheTTL,
		Audience:        audience,
		CustomClaimPath: customClaimPath,
	}
}
