package dialer

import (
	"time"

	"github.com/consensysquorum/quorum-key-manager/pkg/json"
)

type Config struct {
	Timeout   *json.Duration `json:"timeout,omitempty" description:"Max time to wait for a connection to complete (if zero, no timeout)"`
	KeepAlive *json.Duration `json:"keepAlive,omitempty" description:"Interval between keep-alive probes for an active network connection (if zero, if default to 15 seconds)"`
}

func (cfg *Config) SetDefault() *Config {
	if cfg.Timeout == nil {
		cfg.Timeout = &json.Duration{Duration: 30 * time.Second}
	}

	if cfg.KeepAlive == nil {
		cfg.KeepAlive = &json.Duration{Duration: 30 * time.Second}
	}

	return cfg
}
