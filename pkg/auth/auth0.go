package auth

import (
	"fmt"
	"net/http"

	json2 "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/pkg/tls/certificate"
)

const (
	Auth0IssuerServerDomain = "auth0.com"
)

type JWKsResponse struct {
	Keys []JWKsKey `json:"keys"`
}

type JWKsKey struct {
	Alg string   `json:"alg"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	E   string   `json:"e,omitempty"`
	N   string   `json:"n,omitempty"`
	Kid string   `json:"kid"`
	X5c []string `json:"x5c"`
	X5t string   `json:"x5t,omitempty"`
}

func JWKsCertificates(client *http.Client, auth0Domain string) ([]certificate.KeyPair, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/.well-known/jwks.json", auth0Domain), nil)

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp := new(JWKsResponse)
	if err = json2.UnmarshalBody(r.Body, resp); err != nil {
		return nil, err
	}

	keyPairs := []certificate.KeyPair{}
	for _, k := range resp.Keys {
		keyPairs = append(keyPairs, certificate.KeyPair{
			Cert: []byte(k.X5c[0]),
		})
	}

	return keyPairs, nil
}
