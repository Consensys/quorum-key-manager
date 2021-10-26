package app

import (
	"github.com/consensys/quorum-key-manager/src/infra/api-key/csv"
	"github.com/consensys/quorum-key-manager/src/infra/jwt/jose"
	manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
)

type Config struct {
	Manifest *manifests.Config
	OIDC     *jose.Config
	APIKey   *csv.Config
	TLS      interface{}
}
