package auth

import (
	apikey "github.com/consensys/quorum-key-manager/src/auth/authenticator/api-key"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/tls"
	"github.com/consensys/quorum-key-manager/src/infra/http/middlewares/jwt"
	manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
)

type Config struct {
	OIDC     *jwt.Config
	TLS      *tls.Config
	APIKEY   *apikey.Config
	Manifest *manifests.Config
}
