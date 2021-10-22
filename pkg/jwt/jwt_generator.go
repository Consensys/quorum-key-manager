package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenGenerator struct {
	privateKey *rsa.PrivateKey
}

func NewTokenGenerator(key *rsa.PrivateKey) *TokenGenerator {
	return &TokenGenerator{privateKey: key}
}

func (j *TokenGenerator) GenerateAccessToken(claims map[string]interface{}, ttl time.Duration) (tokenValue string, err error) {
	sc := jwt.StandardClaims{
		Issuer:    "quorum-key-manager",
		IssuedAt:  time.Now().UTC().Unix(),
		NotBefore: time.Now().UTC().Unix(),
		Subject:   "",
		ExpiresAt: time.Now().UTC().Add(ttl).Unix(),
	}

	bsc, _ := json.Marshal(sc)

	c := jwt.MapClaims{}
	_ = json.Unmarshal(bsc, &c)

	for k, v := range claims {
		c[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	s, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return s, nil
}
