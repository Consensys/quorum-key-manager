package app

import (
	manifests "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
	pg "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
)

type Config struct {
	Postgres *pg.Config
	Manifest *manifests.Config
}
