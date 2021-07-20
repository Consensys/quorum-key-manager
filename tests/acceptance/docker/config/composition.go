package config

import (
	"github.com/consensys/quorum-key-manager/tests/acceptance/docker/config/hashicorp"
	"github.com/consensys/quorum-key-manager/tests/acceptance/docker/config/postgres"
	"github.com/consensys/quorum-key-manager/tests/acceptance/utils"
)

type Composition struct {
	Containers map[string]*Container
}

type Container struct {
	HashicorpVault *hashicorp.Config
	Postgres       *postgres.Config
}

func (c *Container) Field() (interface{}, error) {
	return utils.ExtractField(c)
}
