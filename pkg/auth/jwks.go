package auth

import (
	"context"
	"fmt"
	"net/http"

	json2 "github.com/consensys/quorum-key-manager/pkg/json"
	"gopkg.in/square/go-jose.v2"
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

func RetrieveKeySet(ctx context.Context, client *http.Client, authEndpoint string) (*jose.JSONWebKeySet, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", authEndpoint, nil)
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call to JWK server failed. %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve keys from JWK server.")
	}

	keySet := &jose.JSONWebKeySet{}
	if err := json2.UnmarshalBody(response.Body, keySet); err != nil {
		return nil, fmt.Errorf("failed to decode response body. %s", err.Error())
	}

	return keySet, nil
}
