package tls

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const AuthMode = "Tls"

type Authenticator struct {
}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	auth := &Authenticator{}

	return auth, nil
}

// Authenticate checks certs and retrieve user Info
// CN field -> Username
// Organization -> Groups
func (authenticator Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract Certificate info from request if any
	// let go without error when no cert found
	if len(req.TLS.PeerCertificates) == 0 {
		return nil, nil
	}
	// first array element is the leaf
	clientCert := *req.TLS.PeerCertificates[0]
	// UserInfo returned is retrieved from cert contents
	return &types.UserInfo{
		Username: clientCert.Subject.CommonName,
		Groups:   clientCert.Subject.Organization,
		AuthMode: AuthMode,
	}, nil
}
