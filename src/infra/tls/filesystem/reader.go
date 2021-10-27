package filesystem

import (
	"context"
	"crypto/x509"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/tls"
	"io/ioutil"
	"os"
)

type Reader struct {
	fs os.FileInfo
}

var _ tls.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	fs, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, errors.ConfigError(err.Error())
	}

	return &Reader{fs: fs}, nil
}

func (r *Reader) Load(_ context.Context) (*x509.CertPool, error) {
	fileContent, err := ioutil.ReadFile(r.fs.Name())
	if err != nil {
		return nil, errors.ConfigError(err.Error())
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(fileContent)
	if !ok {
		return nil, errors.ConfigError("failed to append cert to pool")
	}

	return caCertPool, nil
}
