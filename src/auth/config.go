package auth

import (
	apikey "github.com/consensys/quorum-key-manager/src/auth/authenticator/api-key"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oidc"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/tls"
	manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/reader"
)

type Config struct {
	OIDC     *oidc.Config
	TLS      *tls.Config
	APIKEY   *apikey.Config
	Manifest *manifests.Config
}
