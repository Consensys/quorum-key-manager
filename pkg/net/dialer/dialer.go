package dialer

import "net"

func New(cfg *Config) *net.Dialer {
	return &net.Dialer{
		Timeout:   cfg.Timeout.Duration,
		KeepAlive: cfg.KeepAlive.Duration,
	}
}
