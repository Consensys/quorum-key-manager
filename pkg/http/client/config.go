package httpclient

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
)

// Config for creating an HTTP Client
type Config struct {
	Transport *transport.Config `json:"transport,omitempty"`
	Timeout   *json.Duration    `json:"timeout,omitempty"`
}

func (cfg *Config) SetDefault() *Config {
	if cfg.Transport == nil {
		cfg.Transport = new(transport.Config)
	}
	cfg.Transport.SetDefault()

	if cfg.Timeout == nil {
		cfg.Timeout = &json.Duration{Duration: 0}
	}

	return cfg
}
