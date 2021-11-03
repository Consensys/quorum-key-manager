package jose

import (
	"time"
)

type Config struct {
	IssuerURL string
	CacheTTL  time.Duration
	Audience  []string
}

func NewConfig(issuerURL string, audience []string, cacheTTL time.Duration) *Config {
	return &Config{
		IssuerURL: issuerURL,
		CacheTTL:  cacheTTL,
		Audience:  audience,
	}
}
