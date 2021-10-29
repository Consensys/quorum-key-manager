package client

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"

	"github.com/go-pg/pg/v10"

	"github.com/consensys/quorum-key-manager/pkg/tls"
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
	ApplicationName   string        `json:"application_name"`
	SSLMode           string        `json:"ssl_mode"`
	TLSCert           string        `json:"tls_cert"`
	TLSKey            string        `json:"tls_key"`
	TLSCA             string        `json:"tls_ca"`
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

	tlsOption, err := cfg.getTLSOption()
	if err != nil {
		return nil, err
	}

	dialer, err := NewTLSDialer(cfg.SSLMode, cfg.Host, cfg.KeepAliveInterval, cfg.DialTimeout, tlsOption)
	if err != nil {
		return nil, err
	}

	if dialer != nil {
		opt.Dialer = dialer.DialContext
	} else {
		opt.Dialer = Dialer(cfg.KeepAliveInterval, cfg.DialTimeout).DialContext
	}

	return opt, nil
}

func (cfg *Config) getTLSOption() (*tls.Option, error) {
	tlsOption := &tls.Option{}
	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		cert, err := ioutil.ReadFile(cfg.TLSCert)
		if err != nil {
			return nil, err
		}

		key, err := ioutil.ReadFile(cfg.TLSKey)
		if err != nil {
			return nil, err
		}

		tlsOption.Certificates = []*certificate.KeyPair{{Cert: cert, Key: key}}

		if cfg.TLSCA != "" {
			ca, err := ioutil.ReadFile(cfg.TLSCA)
			if err != nil {
				return nil, err
			}

			tlsOption.CAs = [][]byte{ca}
		}
	}

	return tlsOption, nil
}
