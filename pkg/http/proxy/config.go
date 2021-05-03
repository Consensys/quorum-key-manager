package proxy

import (
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
)

type Config struct {
	Transport     *transport.Config     `json:"transport,omitempty"`
	FlushInterval *json.Duration        `json:"flushInterval,omitempty"`
	Request       *request.ProxyConfig  `json:"request,omitempty"`
	Response      *response.ProxyConfig `json:"response,omitempty"`
}

func (cfg *Config) SetDefault() *Config {
	if cfg == nil {
		cfg = new(Config)
	}

	if cfg.Transport == nil {
		cfg.Transport = new(transport.Config)
	}

	cfg.Transport.SetDefault()

	if cfg.FlushInterval == nil {
		cfg.FlushInterval = &json.Duration{Duration: 100 * time.Millisecond}
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
