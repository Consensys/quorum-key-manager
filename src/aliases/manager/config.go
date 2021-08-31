package aliasmanager

import (
	pg "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
)

type Config struct {
	Postgres *pg.Config
}
