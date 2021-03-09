package config

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/docker/container/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/acceptance-tests/utils"
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
