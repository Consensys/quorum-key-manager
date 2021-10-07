package app

import (
	pg "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
)

// Config store aliases config.
type Config struct {
	Postgres *pg.Config
}
