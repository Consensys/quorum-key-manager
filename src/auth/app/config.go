package app

import (
	apikey "github.com/consensys/quorum-key-manager/src/infra/api-key/csv"
	jwt "github.com/consensys/quorum-key-manager/src/infra/jwt/jose"
	manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/yaml"
	tls "github.com/consensys/quorum-key-manager/src/infra/tls/filesystem"
)

type Config struct {
	Manifest *manifests.Config
	OIDC     *jwt.Config
	APIKey   *apikey.Config
	TLS      *tls.Config
}
