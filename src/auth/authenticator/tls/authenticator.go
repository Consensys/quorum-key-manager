package tls

import (
	"bytes"
	"crypto/x509"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const AuthMode = "Tls"

type Authenticator struct {
	Certificates []*x509.Certificate
}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if len(cfg.Certificates) == 0 {
		return nil, nil
	}

	auth := &Authenticator{Certificates: cfg.Certificates}

	return auth, nil
}

// Authenticate checks certs and retrieve user Info
// CN field -> Username
// Organization -> Groups
func (authenticator Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract Certificate info from request if any
	if len(req.TLS.PeerCertificates) == 0 {
		return types.AnonymousUser, nil
	}
	// first array element is the leaf
	clientCert := *req.TLS.PeerCertificates[0]
	// check this cert matches authenticator provided one
	// using strict comparison
	var matchingCert bool

	for _, authCert := range authenticator.Certificates {
		if bytes.Equal(clientCert.Raw, authCert.Raw) {
			matchingCert = true
			break
		}
	}
	if matchingCert {
		return &types.UserInfo{
			Username: clientCert.Subject.CommonName,
			Groups:   clientCert.Subject.Organization,
			AuthMode: AuthMode,
		}, nil
	}
	return nil, errors.UnauthorizedError("no matching cert found")
}
