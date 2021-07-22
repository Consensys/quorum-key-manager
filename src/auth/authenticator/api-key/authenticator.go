package apikey

import (
	"encoding/base64"
	"encoding/hex"
	"hash"
	"net/http"
	"strings"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const (
	AuthMode    = "ApiKey"
	BasicSchema = "Basic"
)

type Authenticator struct {
	APIKeyFile map[string]*UserNameAndGroups
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
// ? -> Username
// ? -> Groups
func (authenticator Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract ApiKey
	clientAPIKey, ok, err := extractAPIKey(req.Header.Get("Authorization"), authenticator.B64Encoder)
	if err != nil {
		// could not be decoded
		return nil, errors.UnauthorizedError("apikey format error")
	}
	// respond nothing when none found
	if !ok || clientAPIKey == "" {
		return nil, nil
	}

	h := *authenticator.Hasher
	h.Reset()
	h.Write([]byte(clientAPIKey))
	clientAPIKeyHash := h.Sum(nil)

	// search hex string hashes
	strClientHash := hex.EncodeToString(clientAPIKeyHash[:])
	// Upper case hash
	strClientHash = strings.ToUpper(strClientHash)
	userAndGroups, contains := authenticator.APIKeyFile[strClientHash]
	if contains {
		return &types.UserInfo{
			AuthMode: AuthMode,
			Username: userAndGroups.UserName,
			Groups:   userAndGroups.Groups,
		}, nil
	}

	return nil, errors.UnauthorizedError("apikey does not match")
}

func extractAPIKey(auth string, b64encoder *base64.Encoding) (apiKey string, found bool, err error) {
	if len(auth) < len(BasicSchema) || !strings.EqualFold(auth[:len(BasicSchema)], BasicSchema) {
		return "", false, nil
	}
	b64EncodedAPIKey := auth[len(BasicSchema)+1:]
	decodedAPIKey, err := b64encoder.DecodeString(b64EncodedAPIKey)
	if err != nil {
		return "", true, err
	}
	return string(decodedAPIKey), true, nil
}
