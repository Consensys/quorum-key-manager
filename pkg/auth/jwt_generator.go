package auth

import (
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"strings"
	"time"

	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oidc"
	"github.com/golang-jwt/jwt"
)

type JWTGenerator struct {
	privateKey *rsa.PrivateKey
	claims     *oidc.ClaimsConfig
}

func NewJWTGenerator(key crypto.PrivateKey, claims *oidc.ClaimsConfig) (*JWTGenerator, error) {
	return &JWTGenerator{
		privateKey: key.(*rsa.PrivateKey),
		claims:     claims,
	}, nil
}

func (j *JWTGenerator) GenerateAccessToken(username string, groups []string, ttl time.Duration) (tokenValue string, err error) {
	sc := jwt.StandardClaims{
		Issuer:    "quorum-key-manager",
		IssuedAt:  time.Now().UTC().Unix(),
		NotBefore: time.Now().UTC().Unix(),
		Subject:   "test-token",
		ExpiresAt: time.Now().UTC().Add(ttl).Unix(),
	}
	
	bsc, _ := json.Marshal(sc)
	
	c := jwt.MapClaims{}
	_ = json.Unmarshal(bsc, &c)
	
	c[j.claims.Username] = username
	c[j.claims.Group] = strings.Join(groups, ",")

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	s, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return s, nil
}
