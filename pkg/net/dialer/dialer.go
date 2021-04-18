package dialer

import "net"

func New(cfg *Config) *net.Dialer {
	cfg = cfg.Copy().SetDefault()

	return &net.Dialer{
		Timeout:   cfg.Timeout.Duration,
		KeepAlive: cfg.KeepAlive.Duration,
	}
}
