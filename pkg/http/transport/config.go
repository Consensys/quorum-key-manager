package transport

import (
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/net/dialer"
)

//  Config options to configure communication between Traefik and the servers
type Config struct {
	Dialer                *dialer.Config `json:"dialer,omitempty"`
	IdleConnTimeout       *json.Duration `json:"idleConnTimeout,omitempty" description:"Maximum time an idle (keep-alive) connection will remain idle before closing itself (if zero then no limit)"`
	ResponseHeaderTimeout *json.Duration `json:"responseHeaderTimeout,omitempty" description:"Time to wait for a server's response headers after fully writing the request (if zero then no limit)"`
	ExpectContinueTimeout *json.Duration `json:"expectContinueTimeout,omitempty" description:"Time to wait for a server's first response headers after fully writing the request headers if the request has an 'Expect: 100-continue' header (if zero then no limit)"`
	MaxIdleConnsPerHost   int            `json:"maxIdleConnsPerHost,omitempty" description:"Controls the maximum idle (keep-alive) to keep per-host (if zero, defaults to 2)"`
	MaxConnsPerHost       int            `json:"maxConnsPerHost,omitempty" description:"Limits the total number of connections per host, including connections in the dialing, active, and idle states (if zero then no limit)"`
	DisableKeepAlives     bool           `json:"disableKeepAlives,omitempty" description:"Disables HTTP keep-alives and will only use the connection to the server for a single HTTP request"`
	DisableCompression    bool           `json:"disableCompression,omitempty" description:"Prevents from requesting compression with an 'Accept-Encoding: gzip' request header when the Request contains no existing Accept-Encoding value"`
	EnableHTTP2           bool           `json:"enableHTTP2,omitempty" description:"Enables HTTP2 connection"`
	EnableH2C             bool           `json:"enableH2C,omitempty" description:"Enables H2C connection"`
}

func (cfg *Config) SetDefault() {
	if cfg.Dialer == nil {
		cfg.Dialer = &dialer.Config{}
	}

	cfg.Dialer.SetDefault()

	if cfg.IdleConnTimeout == nil {
		cfg.IdleConnTimeout = &json.Duration{Duration: 90 * time.Second}
	}

	if cfg.ResponseHeaderTimeout == nil {
		cfg.ResponseHeaderTimeout = &json.Duration{Duration: 0}
	}

	if cfg.ExpectContinueTimeout == nil {
		cfg.ExpectContinueTimeout = &json.Duration{Duration: time.Second}
	}

}
