package tls

import (
	"crypto/x509"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/utils"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const AuthMode = "Tls"

type Authenticator struct {
	certs []*x509.Certificate
}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if cfg.CAs == nil {
		return nil, nil
	}
	auth := &Authenticator{
		certs: cfg.CAs,
	}
	return auth, nil
}

// Authenticate checks certs and retrieve user Info
func (auth Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract Certificate info from request if any
	// let go without error when no cert found
	if req.TLS == nil || req.TLS.PeerCertificates == nil || len(req.TLS.PeerCertificates) == 0 {
		return nil, nil
	}

	// first array element is the leaf
	clientCert := req.TLS.PeerCertificates[0]

	isAllowed := false
	for _, cert := range auth.certs {
		if cert.Equal(clientCert) {
			isAllowed = true
		}
	}

	if !isAllowed {
		return nil, errors.UnauthorizedError("request certificate is not valid")
	}

	// UserInfo returned is retrieved from cert contents
	userInfo := &types.UserInfo{
		AuthMode: AuthMode,
	}
	userInfo.Username, userInfo.Tenant = utils.ExtractUsernameAndTenant(clientCert.Subject.CommonName)
	userInfo.Roles, userInfo.Permissions = utils.ExtractRolesAndPermission(clientCert.Subject.Organization)

	return userInfo, nil
}
