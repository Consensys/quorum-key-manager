package tls

import (
	"net/http"
	"strings"

	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const AuthMode = "Tls"

type Authenticator struct {
}

var _ authenticator.Authenticator = Authenticator{}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if len(cfg.Certificates) == 0 {
		return nil, nil
	}

	auth := &Authenticator{}

	return auth, nil
}

func (a Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract Certificate info from request if any
	if len(req.TLS.PeerCertificates) > 0 {
		// As mentioned in doc first array element is the leaf
		cert := req.TLS.PeerCertificates[0]
		//TODO check this cert matches a provided one
		return &types.UserInfo{
			Username: cert.Subject.CommonName,
			Groups:   cert.Subject.Organization,
			AuthMode: AuthMode,
		}, nil

	}

	return nil, nil
}

func extractToken(prefix, auth string) (string, bool) {
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}

	return auth[len(prefix)+1:], true
}
