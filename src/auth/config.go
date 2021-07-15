package auth

import (
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oidc"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/tls"
)

type Config struct {
	OIDC *oidc.Config
	// TODO: list config for each authenticator
	TlsCfg *tls.Config
}
