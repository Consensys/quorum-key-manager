package filesystem

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/infra/tls"
	"io/ioutil"
	"os"
)

type Reader struct {
	path string
}

var _ tls.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	_, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}

	return &Reader{path: cfg.Path}, nil
}

func (r *Reader) Load(_ context.Context) (*x509.CertPool, error) {
	fileContent, err := ioutil.ReadFile(r.path)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(fileContent)
	if !ok {
		return nil, fmt.Errorf("failed to append cert to pool")
	}

	return caCertPool, nil
}
