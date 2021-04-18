package httpclient

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
)

// Config for creating an HTTP Client
type Config struct {
	Transport *transport.Config     `json:"transport,omitempty"`
	Timeout   *json.Duration        `json:"timeout,omitempty"`
	Request   *request.ProxyConfig  `json:"request,omitempty"`
	Response  *response.ProxyConfig `json:"response,omitempty"`
}

func (cfg *Config) Copy() *Config {
	if cfg == nil {
		return nil
	}

	return &Config{
		Transport: cfg.Transport.Copy(),
		Timeout:   cfg.Timeout.Copy(),
		Request:   cfg.Request.Copy(),
		Response:  cfg.Response.Copy(),
	}
}

func (cfg *Config) SetDefault() *Config {
	if cfg == nil {
		cfg = new(Config)
	}

	if cfg.Transport == nil {
		cfg.Transport = new(transport.Config)
	}

	cfg.Transport.SetDefault()

	if cfg.Timeout == nil {
		cfg.Timeout = &json.Duration{Duration: 0}
	}

	if cfg.Request == nil {
		cfg.Request = new(request.ProxyConfig)
	}

	cfg.Request.SetDefault()

	if cfg.Response == nil {
		cfg.Response = new(response.ProxyConfig)
	}

	return cfg
}
