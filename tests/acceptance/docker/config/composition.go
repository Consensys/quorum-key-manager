package config

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/tests/acceptance/docker/container/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/tests/acceptance/utils"
)

type Composition struct {
	Containers map[string]*Container
}

type Container struct {
	HashicorpVault *hashicorp.Config
}

func (c *Container) Field() (interface{}, error) {
	return utils.ExtractField(c)
}
