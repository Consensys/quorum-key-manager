package transport

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/net/dialer"
	"golang.org/x/net/http2"
)

// New creates an http.RoundTripper configured with the Transport configuration settings.
func New(cfg *Config) (http.RoundTripper, error) {
	if cfg == nil {
		cfg = new(Config)
	}

	cfg.SetDefault()

	// Create dialer
	dlr := dialer.New(cfg.Dialer)

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dlr.DialContext,
		DisableKeepAlives:     cfg.DisableKeepAlives,
		DisableCompression:    cfg.DisableCompression,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		MaxConnsPerHost:       cfg.MaxConnsPerHost,
		IdleConnTimeout:       cfg.IdleConnTimeout.Duration,
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout.Duration,
		ExpectContinueTimeout: cfg.ExpectContinueTimeout.Duration,
	}

	if cfg.EnableHTTP2 {
		err := http2.ConfigureTransport(transport)
		if err != nil {
			return nil, err
		}
	}

	if cfg.EnableH2C {
		transport.RegisterProtocol("h2c", &h2cTransportWrapper{
			Transport: &http2.Transport{
				DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(netw, addr)
				},
				AllowHTTP:          true,
				DisableCompression: cfg.DisableCompression,
			},
		})
	}

	return transport, nil
}

type h2cTransportWrapper struct {
	*http2.Transport
}

func (t *h2cTransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	return t.Transport.RoundTrip(req)
}
