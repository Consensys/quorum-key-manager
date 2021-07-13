package auth

import (
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oidc"
)

type Config struct {
	OIDC *oidc.Config
}
