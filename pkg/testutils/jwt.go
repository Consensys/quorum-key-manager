package testutils

import (
	"crypto/rsa"
	"encoding/json"
	"time"

	"gopkg.in/square/go-jose.v2/jwt"
)

func GenerateAccessToken(key *rsa.PrivateKey, issuer, aud, subject string, claims map[string]interface{}, ttl time.Duration) (tokenValue string, err error) {
	sc := jwt.Claims{
		ID:       "id",
		Issuer:   issuer,
		Subject:  subject,
		Audience: []string{aud},
	}

	bsc, _ := json.Marshal(sc)

	c := map[string]interface{}{}
	_ = json.Unmarshal(bsc, &c)

	for k, v := range claims {
		c[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	s, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return s, nil
}
