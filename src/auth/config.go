package auth

import (
	apikey "github.com/consensys/quorum-key-manager/src/auth/authenticator/api-key"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oidc"
)

type Config struct {
	OIDC   *oidc.Config
	APIKEY *apikey.Config
}
