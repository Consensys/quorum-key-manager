package proxy

import (
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
)

type Config struct {
	Transport      *transport.Config
	PassHostHeader *bool
	FlushInterval  *json.Duration
}

func (cfg *Config) SetDefault() {
	if cfg.Transport == nil {
		cfg.Transport = new(transport.Config)
	}

	cfg.Transport.SetDefault()

	if cfg.PassHostHeader != nil {
		cfg.PassHostHeader = new(bool)
		*cfg.PassHostHeader = true
	}

	if cfg.FlushInterval == nil {
		cfg.FlushInterval = &json.Duration{Duration: 100 * time.Millisecond}
	}
}
