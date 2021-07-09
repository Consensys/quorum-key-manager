package oicd

type Config struct {
	Certificate string
	Claims      *ClaimsConfig
}

type ClaimsConfig struct {
	Username string
	Group    string
}

func NewConfig(ca, usernameClaim, groupClaims string) *Config {
	return &Config{
		Certificate: ca,
		Claims: &ClaimsConfig{
			Username: usernameClaim,
			Group:    groupClaims,
		},
	}
}
