package tls

import (
	"crypto/x509"
	"github.com/consensys/quorum-key-manager/src/infra/http/middlewares/utils"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/tls"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const AuthMode = "Tls"

type Authenticator struct {
	rootCAs *x509.CertPool
}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if cfg.CAs == nil {
		return nil, nil
	}
	auth := &Authenticator{
		rootCAs: cfg.CAs,
	}
	return auth, nil
}

// Authenticate checks rootCAs and retrieve user Info
func (auth Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract Certificate info from request if any
	// let go without error when no cert found
	if req.TLS == nil || req.TLS.PeerCertificates == nil || len(req.TLS.PeerCertificates) == 0 {
		return nil, nil
	}

	if !req.TLS.HandshakeComplete {
		return nil, errors.UnauthorizedError("request must complete valid handshake")
	}

	err := tls.VerifyCertificateAuthority(req.TLS.PeerCertificates, req.TLS.ServerName, auth.rootCAs, true)
	if err != nil {
		return nil, errors.UnauthorizedError(err.Error())
	}

	// UserInfo returned is retrieved from cert contents
	userInfo := &types.UserInfo{
		AuthMode: AuthMode,
	}

	// first array element is the leaf
	clientCert := req.TLS.PeerCertificates[0]
	userInfo.Username, userInfo.Tenant = utils.ExtractUsernameAndTenant(clientCert.Subject.CommonName)
	userInfo.Permissions = utils.ExtractPermissions(clientCert.Subject.OrganizationalUnit)
	userInfo.Roles = clientCert.Subject.Organization

	return userInfo, nil
}
