package client

import (
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
)

// Config object that be converted into an api.Config later
type Config struct {
	Address       string
	CACert        string
	CAPath        string
	ClientCert    string
	ClientKey     string
	TLSServerName string
	Namespace     string
	ClientTimeout time.Duration
	RateLimit     float64
	BurstLimit    int
	MaxRetries    int
	SkipVerify    bool
	Token         string
}

func NewConfig(addr, namespace string) *Config {
	return &Config{
		Address:   addr,
		Namespace: namespace,
	}
}

// ToHashicorpConfig extracts an api.Config object from self
func (c *Config) ToHashicorpConfig() *api.Config {
	// Create Hashicorp Configuration
	config := api.DefaultConfig()
	config.Address = c.Address
	config.HttpClient = cleanhttp.DefaultClient()
	config.HttpClient.Timeout = time.Second * 60

	return config
}
