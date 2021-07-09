package client

import (
	"fmt"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/tls"
	"github.com/go-pg/pg/v10"
)

const (
	requireSSLMode    = "require"
	disableSSLMode    = "disable"
	verifyCASSLMode   = "verify-ca"
	verifyFullSSLMode = "verify-full"
)

type Config struct {
	Host              string
	Port              string
	User              string
	Password          string
	Database          string
	PoolSize          int
	PoolTimeout       time.Duration
	DialTimeout       time.Duration
	KeepAliveInterval time.Duration
	TLS               *tls.Option
	ApplicationName   string
	SSLMode           string
}

func (cfg *Config) ToPGOptions() (*pg.Options, error) {
	opt := &pg.Options{
		Addr:            fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		User:            cfg.User,
		Password:        cfg.Password,
		Database:        cfg.Database,
		PoolSize:        cfg.PoolSize,
		ApplicationName: cfg.ApplicationName,
		PoolTimeout:     cfg.PoolTimeout,
	}

	dialer, err := NewTLSDialer(cfg)
	if err != nil {
		return nil, err
	}

	if dialer != nil {
		opt.Dialer = dialer.DialContext
	} else {
		opt.Dialer = Dialer(cfg).DialContext
	}

	return opt, nil
}
