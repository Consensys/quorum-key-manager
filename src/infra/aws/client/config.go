package client

type Config struct {
	Endpoint  string
	Region    string
	AccessID  string
	SecretKey string
}

func NewBaseConfig(region, accessID, secretKey string) *Config {
	return &Config{
		Region:    region,
		AccessID:  accessID,
		SecretKey: secretKey,
	}
}

func NewIntegrationConfig(region, endpoint string) *Config {
	return &Config{
		Region:   region,
		Endpoint: endpoint,
	}
}
