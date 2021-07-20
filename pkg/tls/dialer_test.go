package tls

import (
	"context"
	"crypto/tls"
	"net"
	"testing"
)

type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(b []byte) (int, error) { return f(b) }

func TestDialer(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		ln, err = net.Listen("tcp6", "[::1]:0")
	}
	if err != nil {
		t.Fatal(err)
	}

	defer ln.Close()

	unblockServer := make(chan struct{}) // close-only
	defer close(unblockServer)
	go func() {
		conn, e := ln.Accept()
		if e != nil {
			return
		}
		defer conn.Close()
		<-unblockServer
	}()

	ctx, cancel := context.WithCancel(context.Background())
	d := Dialer{
		Dialer: &net.Dialer{},
		TLSConfig: &tls.Config{
			Rand: readerFunc(func(b []byte) (n int, err error) {
				// By the time crypto/tls wants randomness, that means it has a TCP
				// connection, so we're past the Dialer's dial and now blocked
				// in a handshake. Cancel our context and see if we get unstuck.
				// (Our TCP listener above never reads or writes, so the Handshake
				// would otherwise be stuck forever)
				cancel()
				return len(b), nil
			}),
			ServerName: "foo",
		},
	}

	_, err = d.DialContext(ctx, "tcp", ln.Addr().String())
	if err != context.Canceled {
		t.Errorf("err = %v; want context.Canceled", err)
	}
}
