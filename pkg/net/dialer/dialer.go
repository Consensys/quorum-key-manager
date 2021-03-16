package dialer

import "net"

func New(cfg *Config) *net.Dialer {
	if cfg == nil {
		cfg = new(Config)
	}

	cfg.SetDefault()

	return &net.Dialer{
		Timeout:   cfg.Timeout.Duration,
		KeepAlive: cfg.KeepAlive.Duration,
	}
}
