package tls

import (
	"context"
	"crypto/x509"
)

//go:generate mockgen -source=reader.go -destination=mock/reader.go -package=mock

// Reader reads TLS certificates
type Reader interface {
	Load(ctx context.Context) (*x509.CertPool, error)
}
