package tls

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"net/http"
)

const AuthMode = "Tls"

type Authenticator struct {
	Certificates []*x509.Certificate
}

var _ authenticator.Authenticator = Authenticator{}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if len(cfg.Certificates) == 0 {
		return nil, nil
	}

	auth := &Authenticator{Certificates: cfg.Certificates}

	return auth, nil
}

func (authenticator Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract Certificate info from request if any
	if len(req.TLS.PeerCertificates) > 0 {
		// As mentioned in doc first array element is the leaf
		cert := req.TLS.PeerCertificates[0]
		//check this cert matches authenticator provided one
		// using strict comparison
		var matchingCert bool

		for i, cert := range authenticator.Certificates {
			if bytes.Equal(cert.Raw, authenticator.Certificates[i].Raw) {
				matchingCert = true
				break
			}
		}
		if matchingCert {
			return &types.UserInfo{
				Username: cert.Subject.CommonName,
				Groups:   cert.Subject.Organization,
				AuthMode: AuthMode,
			}, nil
		} else {
			return nil, fmt.Errorf("certs do not match")
		}

	}

	return nil, fmt.Errorf("request provides no cert")
}
