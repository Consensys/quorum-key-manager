package oidc

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type JWTChecker struct {
	certs     []*x509.Certificate
	parser    *jwt.Parser
	claimsCfg *ClaimsConfig
}

func NewJWTChecker(certs []*x509.Certificate, claimsCfg *ClaimsConfig, skipClaimsValidation bool) *JWTChecker {
	return &JWTChecker{
		certs:     certs,
		claimsCfg: claimsCfg,
		parser: &jwt.Parser{
			SkipClaimsValidation: skipClaimsValidation,
		},
	}
}

func (checker *JWTChecker) Check(_ context.Context, bearerToken string) (*Claims, error) {
	if len(checker.certs) == 0 {
		// If no certificate provided we deactivate authentication
		return nil, nil
	}

	// Parse and validate token injected in context
	token, err := checker.parser.ParseWithClaims(
		bearerToken,
		&Claims{cfg: checker.claimsCfg},
		checker.keyFunc,
	)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return token.Claims.(*Claims), nil
}

func (checker *JWTChecker) keyFunc(token *jwt.Token) (interface{}, error) {
	for _, cert := range checker.certs {
		if pubkey, err := tokenAlgoChecker(token.Method.Alg(), cert); err == nil {
			return pubkey, nil
		}
	}

	return nil, fmt.Errorf("unable to find appropriate key in key set")
}

func tokenAlgoChecker(method string, cert *x509.Certificate) (interface{}, error) {
	switch method {
	case "RS256", "RS384", "RS512":
		pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate is not an RSA public key")
		}
		return pubKey, nil
	case "ES256", "ES384", "ES512":
		pubKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate is not an ECDSA public key")
		}
		return pubKey, nil
	default:
		return nil, fmt.Errorf("invalid access token signing method %q", method)
	}
}
