package oicd

import (
	"net/http"
	"strings"

	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const AuthMode = "JWT"

type Authenticator struct {
	jwtCheckers []*JWTChecker
}

var _ authenticator.Authenticator = Authenticator{}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if len(cfg.Certificates) == 0 {
		return nil, nil
	}

	auth := &Authenticator{
		jwtCheckers: []*JWTChecker{},
	}

	for _, cert := range cfg.Certificates {
		keyPair, err := certificate.X509KeyPair([]byte(cert), nil)
		if err != nil {
			return nil, err
		}

		auth.jwtCheckers = append(auth.jwtCheckers, NewJWTChecker(keyPair.Leaf, cfg.Claims, false))
	}

	return auth, nil
}

func (a Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// Extract Access Token from context
	token, ok := extractToken("Bearer ", req.Header.Get("Authorization"))
	if !ok {
		return nil, nil
	}

	var err error
	var claims *Claims
	for _, jwtChecker := range a.jwtCheckers {
		claims, err = jwtChecker.Check(req.Context(), token)
		if err == nil {
			return &types.UserInfo{
				Username: claims.Username,
				Groups:   claims.Groups,
				AuthMode: AuthMode,
			}, nil
		}
	}
	
	return nil, err
}

func extractToken(prefix, auth string) (string, bool) {
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}

	return auth[len(prefix):], true
}
