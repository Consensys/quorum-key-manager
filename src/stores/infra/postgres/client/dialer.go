package client

import (
	"context"
	gotls "crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/tcp"
	"github.com/consensys/quorum-key-manager/pkg/tls"
)

type SSLDialer struct {
	Dialer tcp.Dialer
}

func Dialer(cfg *Config) *net.Dialer {
	return &net.Dialer{
		Timeout:   cfg.DialTimeout,
		KeepAlive: cfg.KeepAliveInterval,
	}
}

func (d *SSLDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := d.Dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	// Next is a preliminary send/receive message between client and Postgres server
	// to make sure server is configured for TLS connection
	//
	// It should happened before upgrading the connection to TLS
	// Implementation is largely inspired from https://github.com/lib/pq/blob/v1.7.0/conn.go#L1027
	var scratch [512]byte
	scratch[0] = 0
	buf := scratch[:5]

	x := make([]byte, 4)
	binary.BigEndian.PutUint32(x, uint32(80877103))
	buf = append(buf, x...)

	wrap := buf[1:]
	binary.BigEndian.PutUint32(wrap, uint32(len(wrap)))

	_, err = conn.Write(buf[1:])
	if err != nil {
		conn.Close()
		return nil, err
	}

	b := scratch[:1]
	_, err = io.ReadFull(conn, b)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if b[0] != 'S' {
		conn.Close()
		return nil, fmt.Errorf("ssl is not enabled on the server")
	}

	return conn, nil
}

type TLSDialer struct {
	Dialer       *tls.Dialer
	verifyCAOnly bool
}

func NewTLSDialer(cfg *Config) (*TLSDialer, error) {
	var verifyCAOnly bool
	switch cfg.SSLMode {
	case requireSSLMode:
		// Setting InsecureSkipVerify to true
		// makes client skip server certificate verification
		// at handshake
		cfg.TLS.InsecureSkipVerify = true
	case verifyFullSSLMode:
		// Setting ServerName
		// makes client proceed to server certificate verification
		// at handshake
		//
		// In this case it controls both
		// - server certificate is CA signed if CA has been passed
		// - server that is accessed is listed in server certificate domains
		cfg.TLS.ServerName = cfg.Host
	case verifyCASSLMode:
		// golang crypto/tls does not allow to implement
		// verify-ca behavior (only verify-full)
		// so we need some customisation
		cfg.TLS.InsecureSkipVerify = true
		verifyCAOnly = true
	case disableSSLMode, "":
		return nil, nil
	default:
		return nil, errors.ConfigError("invalid sslmode")
	}

	tlsConfig, err := tls.NewConfig(cfg.TLS)
	if err != nil {
		return nil, err
	}

	// In case TLS is activated we use a custom dialer that allows
	// sslmode 'verify-ca'
	return &TLSDialer{
		Dialer: &tls.Dialer{
			Dialer: &SSLDialer{
				Dialer: &net.Dialer{
					Timeout:   cfg.DialTimeout,
					KeepAlive: cfg.KeepAliveInterval,
				},
			},
			TLSConfig: tlsConfig,
		},
		verifyCAOnly: verifyCAOnly,
	}, nil
}

func (d *TLSDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := d.Dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	if d.verifyCAOnly {
		err = tls.VerifyCertificateAuthority(conn.(*gotls.Conn), d.Dialer.TLSConfig)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
