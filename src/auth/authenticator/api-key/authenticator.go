package apikey

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"hash"
	"net/http"
)

const (
	AuthMode     = "ApiKey"
	ApiKeyHeader = "X-Key-Manager-APIKEY"
)

type Authenticator struct {
	ApiKeyFile map[string]*UserNameAndGroups
	Hasher     hash.Hash
}

var _ authenticator.Authenticator = Authenticator{}

func NewAuthenticator(cfg *Config) (*Authenticator, error) {
	if len(cfg.ApiKeyFile) == 0 {
		return nil, nil
	}

	auth := &Authenticator{ApiKeyFile: cfg.ApiKeyFile,
		Hasher: cfg.Hasher}

	return auth, nil
}

// Authenticate checks APIKEY hashes retrieve user Info
// ? -> Username
// ? -> Groups
func (authenticator Authenticator) Authenticate(req *http.Request) (*types.UserInfo, error) {
	// extract ApiKey
	apiKey := req.Header.Get(ApiKeyHeader)

	clientApiKeyHash := authenticator.Hasher.Sum([]byte(apiKey))

	// compare hashes
	userAndGroups, contains := authenticator.ApiKeyFile[string(clientApiKeyHash)]
	if contains {
		return &types.UserInfo{
			AuthMode: AuthMode,
			Username: userAndGroups.UserName,
			Groups:   userAndGroups.Groups,
		}, nil
	}

	return nil, errors.UnauthorizedError("apikey does not match")
}
