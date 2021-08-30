package apikey

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"strings"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/utils"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const (
	AuthMode    = "ApiKey"
	BasicSchema = "Basic"
)

type Authenticator struct {
	APIKeyFile map[string]UserClaims
	Hasher     *hash.Hash
	B64Encoder *base64.Encoding
}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if len(cfg.APIKeyFile) == 0 {
		return nil, nil
	}

	auth := &Authenticator{APIKeyFile: cfg.APIKeyFile,
		Hasher:     cfg.Hasher,
		B64Encoder: cfg.B64Encoder,
	}

	return auth, nil
}

// Authenticate checks APIKEY hashes retrieve user Info
// ? -> Subject
// ? -> Roles
func (authenticator Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// In case of no credentials are sent we authenticate with Anonymous user
	if req.Header.Get("Authorization") == "" {
		return nil, nil
	}

	// extract ApiKey
	clientAPIKey, err := extractAPIKey(req.Header.Get("Authorization"), authenticator.B64Encoder)
	if err != nil {
		// could not be decoded
		return nil, errors.UnauthorizedError(err.Error())
	}

	h := *authenticator.Hasher
	h.Reset()
	_, err = h.Write([]byte(clientAPIKey))
	if err != nil {
		// could not be written
		return nil, errors.UnauthorizedError(err.Error())
	}
	clientAPIKeyHash := h.Sum(nil)

	// search hex string hashes
	strClientHash := hex.EncodeToString(clientAPIKeyHash)

	auth, ok := authenticator.APIKeyFile[strClientHash]
	if !ok {
		return nil, errors.UnauthorizedError("api-key does not match")
	}

	userInfo := &types.UserInfo{
		AuthMode:    AuthMode,
		Roles:       []string{},
		Permissions: []types.Permission{},
	}

	userInfo.Username, userInfo.Tenant = utils.ExtractUsernameAndTenant(auth.UserName)
	userInfo.Permissions = utils.ExtractPermissions(auth.Claims)

	return userInfo, nil
}

func extractAPIKey(auth string, b64encoder *base64.Encoding) (apiKey string, err error) {
	if len(auth) <= len(BasicSchema) || !strings.EqualFold(auth[:len(BasicSchema)], BasicSchema) {
		return "", fmt.Errorf("apikey was not provided")
	}
	b64EncodedAPIKey := auth[len(BasicSchema)+1:]
	decodedAPIKey, err := b64encoder.DecodeString(b64EncodedAPIKey)
	if err != nil {
		return "", fmt.Errorf("apikey encoding is not supported")
	}
	return string(decodedAPIKey), nil
}
