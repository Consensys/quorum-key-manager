package tcp

import (
	"context"
	"net"
)

type Dialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}
