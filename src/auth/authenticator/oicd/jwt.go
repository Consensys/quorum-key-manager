package oicd

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/golang-jwt/jwt"
)

type JWTChecker struct {
	cert      *x509.Certificate
	parser    *jwt.Parser
	claimsCfg *ClaimsConfig
}

func NewJWTChecker(cert *x509.Certificate, claimsCfg *ClaimsConfig, skipClaimsValidation bool) *JWTChecker {
	return &JWTChecker{
		cert: cert,
		parser: &jwt.Parser{
			SkipClaimsValidation: skipClaimsValidation,
		},
	}
}

func (checker *JWTChecker) Check(ctx context.Context, bearerToken string) (*Claims, error) {
	if checker.cert == nil {
		// If no certificate provided we deactivate authentication
		return nil, nil
	}

	// Parse and validate token injected in context
	token, err := checker.parser.ParseWithClaims(
		bearerToken,
		&Claims{cfg: checker.claimsCfg},
		checker.key,
	)
	if err != nil {
		return nil, errors.UnauthorizedError(err.Error())
	}
	if !token.Valid {
		return nil, errors.UnauthorizedError("invalid access token")
	}

	return token.Claims.(*Claims), nil
}

func (checker *JWTChecker) key(token *jwt.Token) (interface{}, error) {
	switch token.Method.Alg() {
	case "RS256", "RS384", "RS512":
		pubKey, ok := checker.cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate is not an RSA public key")
		}
		return pubKey, nil
	case "ES256", "ES384", "ES512":
		pubKey, ok := checker.cert.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate is not an ECDSA public key")
		}
		return pubKey, nil
	default:
		return nil, fmt.Errorf("unsupported token method signature %q", token.Method.Alg())
	}
}
