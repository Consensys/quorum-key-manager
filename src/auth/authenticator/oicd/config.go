package oicd

type Config struct {
	Certificates   []string
	Claims         *ClaimsConfig
}

type ClaimsConfig struct {
	Username string
	Group    string
}

func NewConfig(usernameClaim, groupClaims string, certs ...string) *Config {
	return &Config{
		Certificates:   certs,
		Claims: &ClaimsConfig{
			Username: usernameClaim,
			Group:    groupClaims,
		},
	}
}
