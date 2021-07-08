package oicd

import (
	"net/http"
	"strings"

	"github.com/consensys/quorum-key-manager/pkg/tls"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const AuthMode = "JWT"

type Authenticator struct {
	jwtChecker *JWTChecker
}

var _ authenticator.Authenticator = Authenticator{}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if cfg.Certificate == "" && cfg.CertificateServer == "" {
		return nil, nil
	}

	cert, err := tls.X509KeyPair([]byte(cfg.Certificate), nil)
	if err != nil {
		return nil, err
	}

	return &Authenticator{
		jwtChecker: NewJWTChecker(cert.Leaf, cfg.Claims, false),
	}, nil
}

func (a Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// Extract Access Token from context
	token, ok := extractToken("Bearer ", req.Header.Get("Authorization"))
	if !ok {
		return nil, nil
	}

	claims, err := a.jwtChecker.Check(req.Context(), token)
	if err != nil {
		return nil, err
	}

	return &types.UserInfo{
		Username: claims.Username,
		Groups:   claims.Groups,
		AuthMode: AuthMode,
	}, nil
}

func extractToken(prefix, auth string) (string, bool) {
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}

	return auth[len(prefix):], true
}
