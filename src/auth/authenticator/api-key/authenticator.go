package apikey

import (
	"hash"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const (
	AuthMode     = "ApiKey"
	APIKeyHeader = "X-Key-Manager-APIKEY"
)

type Authenticator struct {
	APIKeyFile map[string]*UserNameAndGroups
	Hasher     hash.Hash
}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if len(cfg.APIKeyFile) == 0 {
		return nil, nil
	}

	auth := &Authenticator{APIKeyFile: cfg.APIKeyFile,
		Hasher: cfg.Hasher}

	return auth, nil
}

// Authenticate checks APIKEY hashes retrieve user Info
// ? -> Username
// ? -> Groups
func (authenticator Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract ApiKey
	apiKey := req.Header.Get(APIKeyHeader)

	clientAPIKeyHash := authenticator.Hasher.Sum([]byte(apiKey))

	// compare hashes
	userAndGroups, contains := authenticator.APIKeyFile[string(clientAPIKeyHash)]
	if contains {
		return &types.UserInfo{
			AuthMode: AuthMode,
			Username: userAndGroups.UserName,
			Groups:   userAndGroups.Groups,
		}, nil
	}

	return nil, errors.UnauthorizedError("apikey does not match")
}
