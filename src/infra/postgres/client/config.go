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
	Host              string        `json:"host"`
	Port              string        `json:"port"`
	User              string        `json:"user"`
	Password          string        `json:"password"`
	Database          string        `json:"database"`
	PoolSize          int           `json:"pool_size"`
	PoolTimeout       time.Duration `json:"pool_timeout"`
	DialTimeout       time.Duration `json:"dial_timeout"`
	KeepAliveInterval time.Duration `json:"keep_alive_interval"`
	TLS               *tls.Option   `json:"tls"`
	ApplicationName   string        `json:"application_name"`
	SSLMode           string        `json:"ssl_mode"`
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
