package tls

import (
	"context"
	"crypto/tls"
	"net"
	"strings"

	"github.com/consensys/quorum-key-manager/pkg/tcp"
)

type Dialer struct {
	Dialer tcp.Dialer

	TLSConfig *tls.Config
}

func (d *Dialer) Dial(network, addr string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, addr)
}

func (d *Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	netConn, err := d.Dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	colonPos := strings.LastIndex(addr, ":")
	if colonPos == -1 {
		colonPos = len(addr)
	}
	hostname := addr[:colonPos]

	// If no ServerName is set, infer the ServerName
	// from the hostname we're connecting to.
	config := d.TLSConfig
	if d.TLSConfig.ServerName == "" {
		// Make a copy to avoid polluting argument or default.
		config = d.TLSConfig.Clone()
		config.ServerName = hostname
	}

	tlsConn := tls.Client(netConn, config)

	handshakeErrors := make(chan error, 1)
	go func() {
		handshakeErrors <- tlsConn.Handshake()
		close(handshakeErrors)
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-handshakeErrors:
		if err == nil {
			return tlsConn, nil
		}
	}

	tlsConn.Close()

	return nil, err
}
